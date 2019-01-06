package contextualosmetalinkrepository

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/metalink"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
	"github.com/dpb587/metalink/template"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ./{os_name}/{os_version}/*.meta4
type index struct {
	name                string
	logger              logrus.FieldLogger
	config              Config
	repository          *repository.Repository
	releaseVersionIndex releaseversiondatastore.Index
}

var _ datastore.Index = &index{}

func New(name string, releaseVersionIndex releaseversiondatastore.Index, config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		name:                name,
		logger:              logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		releaseVersionIndex: releaseVersionIndex,
		config:              config,
		repository:          repository.NewRepository(logger, config.RepositoryConfig),
	}
}

func (i *index) GetName() string {
	return i.name
}

func (i *index) GetCompilationArtifacts(f datastore.FilterParams) ([]compilation.Artifact, error) {
	if !f.Release.NameSatisfied(i.config.Release) {
		return nil, nil
	}

	releases, err := i.releaseVersionIndex.GetArtifacts(f.Release, releaseversiondatastore.SingleArtifactLimitParams)
	if err != nil {
		// TODO warn and continue?
		return nil, errors.Wrap(err, "finding release")
	}

	release := releases[0]

	err = i.repository.Reload()
	if err != nil {
		return nil, errors.Wrap(err, "reloading repository")
	}

	paths, err := filepath.Glob(i.repository.Path(i.config.Path, "*", "*", "*.meta4"))
	if err != nil {
		return nil, errors.Wrap(err, "globbing")
	}

	var results = []compilation.Artifact{}

	for _, compiledReleasePath := range paths {
		compiledReleaseBytes, err := ioutil.ReadFile(compiledReleasePath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", compiledReleasePath, err)
		}

		var compiledReleaseMeta4 metalink.Metalink

		err = metalink.Unmarshal(compiledReleaseBytes, &compiledReleaseMeta4)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("unmarshalling %s: %v", compiledReleasePath, err)
		}

		if compiledReleaseMeta4.Files[0].Version == "" {
			i.logger.Warnf("expected version field in file: %s", compiledReleasePath) // TODO make relative to repository path

			continue
		} else if !f.Release.VersionSatisfied(compiledReleaseMeta4.Files[0].Version) {
			continue
		}

		// TODO sanity checks? files = 1?

		osVersionReference := osversion.Reference{
			Name:    path.Base(path.Dir(path.Dir(compiledReleasePath))),
			Version: path.Base(path.Dir(compiledReleasePath)),
		}

		if !f.OS.NameSatisfied(osVersionReference.Name) {
			continue
		} else if !f.OS.VersionSatisfied(osVersionReference.Version) {
			continue
		}

		results = append(
			results,
			compilation.New(
				i.name,
				compilation.Reference{
					ReleaseVersion: release.Reference().(releaseversion.Reference),
					OSVersion:      osVersionReference,
				},
				compiledReleaseMeta4.Files[0],
			),
		)
	}

	return results, nil
}

func (i *index) StoreCompilationArtifact(artifact compilation.Artifact) error {
	artifactRef := artifact.Reference().(compilation.Reference)

	// TODO keep? the compilation -> releaseVersionIndex affiliation should be managed by the config
	if artifactRef.ReleaseVersion.Name != i.config.Release {
		return datastore.UnsupportedOperationErr
	}

	urlLoader := urldefaultloader.New()

	file := artifact.MetalinkFile()

	local, err := urlLoader.Load(metalink.URL{URL: file.URLs[0].URL})
	if err != nil {
		return errors.Wrap(err, "parsing origin destination")
	}

	mirroredFile := metalink.File{
		Name:    file.Name,
		Size:    file.Size,
		Hashes:  file.Hashes,
		Version: artifactRef.ReleaseVersion.Version,
	}

	for _, mirror := range i.config.StorageConfig {
		tmpl, err := template.New(mirror.URI)
		if err != nil {
			return errors.Wrap(err, "parsing mirror destination")
		}

		mirrorWriterURI, err := tmpl.ExecuteString(mirroredFile)
		if err != nil {
			return errors.Wrap(err, "generating mirror uri")
		}

		for k, v := range mirror.Options {
			// TODO unset/revert after?
			os.Setenv(k, v)
		}

		remote, err := urlLoader.Load(metalink.URL{URL: mirrorWriterURI})
		if err != nil {
			return errors.Wrap(err, "loading upload destination")
		}

		progress := pb.New64(int64(file.Size)).Set(pb.Bytes, true).SetRefreshRate(time.Second).SetWidth(80)
		progress.Start()

		err = remote.WriteFrom(local, progress)
		if err != nil {
			return errors.Wrap(err, "copying blob")
		}

		progress.Finish()

		mirroredFile.URLs = append(mirroredFile.URLs, metalink.URL{
			URL:      remote.ReaderURI(),
			Location: mirror.Location,
			Priority: mirror.Priority,
		})
	}

	commitMeta4 := metalink.Metalink{
		Files:     []metalink.File{mirroredFile},
		Generator: "boshua/contextualosmetalinkrepository",
	}

	commitMeta4Bytes, err := metalink.MarshalXML(commitMeta4)
	if err != nil {
		return errors.Wrap(err, "marshalling metalink")
	}

	path := filepath.Join(
		i.config.Path,
		artifactRef.OSVersion.Name,
		artifactRef.OSVersion.Version,
		fmt.Sprintf("v%s.meta4", artifactRef.ReleaseVersion.Version),
	)

	return i.repository.Commit(
		map[string][]byte{path: commitMeta4Bytes},
		fmt.Sprintf(
			"Compiling v%s for %s/%s",
			artifactRef.ReleaseVersion.Version,
			artifactRef.OSVersion.Name,
			artifactRef.OSVersion.Version,
		),
	)
}

func (i *index) FlushCompilationCache() error {
	// TODO defer reload?
	return i.repository.ForceReload()
}
