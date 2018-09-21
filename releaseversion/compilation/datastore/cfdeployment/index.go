package cfdeployment

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/cppforlife/go-patch/patch"
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
	"github.com/dpb587/boshua/deployment/manifest"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore/inmemory"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type index struct {
	name                string
	logger              logrus.FieldLogger
	config              Config
	repository          *repository.Repository
	releaseVersionIndex releaseversiondatastore.Index

	cache      *inmemory.Index
	cacheMutex *sync.Mutex
	cacheWarm  bool
}

var _ datastore.Index = &index{}

func New(name string, releaseVersionIndex releaseversiondatastore.Index, config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		name:                name,
		logger:              logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		releaseVersionIndex: releaseVersionIndex,
		config:              config,
		repository:          repository.NewRepository(logger, config.RepositoryConfig),
		cache:               inmemory.New(),
		cacheMutex:          &sync.Mutex{},
	}
}

func (i *index) GetName() string {
	return i.name
}

type manifestReleases struct {
	Releases []struct {
		Name     string `yaml:"name"`
		Version  string `yaml:"version"`
		Sha1     string `yaml:"sha1"`
		URL      string `yaml:"url"`
		Stemcell struct {
			OS      string `yaml:"os"`
			Version string `yaml:"version"`
		} `yaml:"stemcell"`
	} `yaml:"releases"`
}

type aggregatedKey struct {
	Release        string
	ReleaseVersion string
	ReleaseSha1    string
	OS             string
	OSVersion      string
}

func (i *index) extractCompilations(opsPath string, commitLoader func(string) (*manifestReleases, patch.Ops, error)) (map[aggregatedKey]manifest.ReleasePatch, error) {
	commitList := bytes.NewBuffer(nil)

	err := i.repository.ExecCapture(commitList, "log", "--format=%H", "--", opsPath)
	if err != nil {
		return nil, errors.Wrap(err, "loading file history")
	}

	aggregatedPatches := map[aggregatedKey]manifest.ReleasePatch{}

	scanner := bufio.NewScanner(commitList)
	for scanner.Scan() {
		err := func() error {
			baseManifest, compilationOps, err := commitLoader(scanner.Text())
			if err != nil {
				return errors.Wrapf(err, "loading commit")
			}

			patches := map[string]manifest.ReleasePatch{}

			for _, baseRelease := range baseManifest.Releases {
				patches[baseRelease.Name] = manifest.ReleasePatch{
					Name:    baseRelease.Name,
					Version: baseRelease.Version,
					Source: manifest.ReleasePatchRef{
						Sha1: baseRelease.Sha1,
						URL:  baseRelease.URL,
					},
				}
			}

			baseYaml, err := yaml.Marshal(baseManifest)
			if err != nil {
				return errors.Wrap(err, "remarshalling")
			}

			var baseRawManifest interface{}

			err = yaml.Unmarshal(baseYaml, &baseRawManifest)
			if err != nil {
				return errors.Wrap(err, "unmarshalling raw")
			}

			compilationRawManifest, err := compilationOps.Apply(baseRawManifest)
			if err != nil {
				return errors.Wrap(err, "applying compilation ops")
			}

			compilationYaml, err := yaml.Marshal(compilationRawManifest)
			if err != nil {
				return errors.Wrap(err, "marshalling compilation manifest")
			}

			var compilationManifest manifestReleases

			err = yaml.Unmarshal(compilationYaml, &compilationManifest)
			if err != nil {
				return errors.Wrap(err, "unmarshalling compilation manifest")
			}

			for _, compilationRelease := range compilationManifest.Releases {
				matchRelease, ok := patches[compilationRelease.Name]
				if !ok {
					continue
				} else if matchRelease.Version != compilationRelease.Version {
					continue
				} else if compilationRelease.Stemcell.OS == "" || compilationRelease.Stemcell.Version == "" {
					continue
				}

				matchRelease.Compiled.Sha1 = compilationRelease.Sha1
				matchRelease.Compiled.URL = compilationRelease.URL
				matchRelease.Stemcell.OS = compilationRelease.Stemcell.OS
				matchRelease.Stemcell.Version = compilationRelease.Stemcell.Version

				aggregatedPatches[aggregatedKey{
					Release:        matchRelease.Name,
					ReleaseVersion: matchRelease.Version,
					ReleaseSha1:    matchRelease.Source.Sha1,
					OS:             matchRelease.Stemcell.OS,
					OSVersion:      matchRelease.Stemcell.Version,
				}] = matchRelease
			}

			return nil
		}()
		if err != nil {
			i.logger.Warnf("%s", errors.Wrapf(err, "extracting commit %s", scanner.Text()))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "scanning")
	}

	return aggregatedPatches, nil
}

func (i *index) GetCompilationArtifacts(f datastore.FilterParams) ([]compilation.Artifact, error) {
	err := i.fillCache()
	if err != nil {
		return nil, err
	}

	return i.cache.GetCompilationArtifacts(f)
}

func (i *index) fillCache() error {
	i.cacheMutex.Lock()
	defer i.cacheMutex.Unlock()

	if i.cacheWarm && i.repository.WarmCache() {
		return nil
	}

	err := i.cache.FlushCompilationCache()
	if err != nil {
		return errors.Wrap(err, "flushing in-memory")
	}

	err = i.repository.Reload()
	if err != nil {
		return errors.Wrap(err, "reloading repository")
	}

	aggregatedPatches, err := i.extractCompilations(
		"operations/use-compiled-releases.yml",
		func(commit string) (*manifestReleases, patch.Ops, error) {
			baseManifest := bytes.NewBuffer(nil)

			err := i.repository.ExecCapture(baseManifest, "show", fmt.Sprintf("%s:cf-deployment.yml", commit))
			if err != nil {
				return nil, nil, errors.Wrapf(err, "reading cf-deployment.yml")
			}

			var dep manifestReleases

			err = yaml.Unmarshal(baseManifest.Bytes(), &dep)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "parsing cf-deployment.yml")
			}

			compBytes := bytes.NewBuffer(nil)

			err = i.repository.ExecCapture(compBytes, "show", fmt.Sprintf("%s:operations/use-compiled-releases.yml", commit))
			if err != nil {
				return nil, nil, errors.Wrap(err, "reading use-compiled-releases.yml")
			}

			var compOpDefs []patch.OpDefinition

			err = yaml.Unmarshal(compBytes.Bytes(), &compOpDefs)
			if err != nil {
				return nil, nil, errors.Wrap(err, "parsing use-compiled-releases.yml")
			}

			compOps, err := patch.NewOpsFromDefinitions(compOpDefs)
			if err != nil {
				return nil, nil, errors.Wrap(err, "building ops")
			}

			return &dep, compOps, nil
		},
	)
	if err != nil {
		return errors.Wrap(err, "extracting default compilations")
	}

	xenialPatches, err := i.extractCompilations(
		"operations/experimental/use-compiled-releases-xenial-stemcell.yml",
		func(commit string) (*manifestReleases, patch.Ops, error) {
			baseManifest := bytes.NewBuffer(nil)

			err := i.repository.ExecCapture(baseManifest, "show", fmt.Sprintf("%s:cf-deployment.yml", commit))
			if err != nil {
				return nil, nil, errors.Wrapf(err, "reading cf-deployment.yml")
			}

			var dep manifestReleases

			err = yaml.Unmarshal(baseManifest.Bytes(), &dep)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "parsing cf-deployment.yml")
			}

			compBytes := bytes.NewBuffer(nil)

			err = i.repository.ExecCapture(compBytes, "show", fmt.Sprintf("%s:operations/experimental/use-compiled-releases-xenial-stemcell.yml", commit))
			if err != nil {
				return nil, nil, errors.Wrap(err, "reading use-compiled-releases-xenial-stemcell.yml")
			}

			var compOpDefs []patch.OpDefinition

			err = yaml.Unmarshal(compBytes.Bytes(), &compOpDefs)
			if err != nil {
				return nil, nil, errors.Wrap(err, "parsing use-compiled-releases-xenial-stemcell.yml")
			}

			compOps, err := patch.NewOpsFromDefinitions(compOpDefs)
			if err != nil {
				return nil, nil, errors.Wrap(err, "building ops")
			}

			return &dep, compOps, nil
		},
	)
	if err != nil {
		return errors.Wrap(err, "extracting xenial compilations")
	}

	for k, v := range xenialPatches {
		aggregatedPatches[k] = v
	}

	for _, patch := range aggregatedPatches {
		release, err := releaseversiondatastore.GetArtifact(i.releaseVersionIndex, releaseversiondatastore.FilterParams{
			NameExpected:     true,
			Name:             patch.Name,
			VersionExpected:  true,
			Version:          patch.Version,
			ChecksumExpected: true,
			Checksum:         fmt.Sprintf("sha1:%s", patch.Source.Sha1), // TODO sha256 support
		})
		if err != nil {
			// TODO warn and continue?
			return errors.Wrapf(err, "finding release %s/%s", patch.Name, patch.Version)
		}

		i.cache.Add(
			compilation.New(
				i.name,
				compilation.Reference{
					ReleaseVersion: release.Reference().(releaseversion.Reference),
					OSVersion: osversion.Reference{
						Name:    patch.Stemcell.OS,
						Version: patch.Stemcell.Version,
					},
				},
				metalink.File{
					Name: filepath.Base(patch.Compiled.URL),
					URLs: []metalink.URL{
						{
							URL: patch.Compiled.URL,
						},
					},
					Hashes: []metalink.Hash{
						{
							Type: "sha-1", // TODO
							Hash: patch.Compiled.Sha1,
						},
					},
				},
			),
		)
	}

	return nil
}

func (i *index) StoreCompilationArtifact(artifact compilation.Artifact) error {
	return datastore.UnsupportedOperationErr
}

func (i *index) FlushCompilationCache() error {
	// TODO defer reload?
	return i.repository.ForceReload()
}
