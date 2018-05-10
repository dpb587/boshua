package boshreleasedpb

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore/inmemory"
	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger     logrus.FieldLogger
	config     Config
	repository *git.Repository
	inmemory   datastore.Index
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	idx := &index{
		logger:     logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config:     config,
		repository: git.NewRepository(logger, config.RepositoryConfig),
	}

	idx.inmemory = inmemory.New(idx.loader, idx.repository.Reload)

	return idx
}

func (i *index) Filter(ref compiledreleaseversion.Reference) ([]compiledreleaseversion.Artifact, error) {
	return i.inmemory.Filter(ref)
}

func (i *index) Find(ref compiledreleaseversion.Reference) (compiledreleaseversion.Artifact, error) {
	return i.inmemory.Find(ref)
}

func (i *index) Store(artifact compiledreleaseversion.Artifact) error {
	path := filepath.Join(
		"compiled-release",
		i.config.Channel,
		artifact.OSVersion.Name,
		artifact.OSVersion.Version,
		fmt.Sprintf("%s.meta4", artifact.ReleaseVersion.Version),
	)

	meta4 := metalink.Metalink{
		Files:     []metalink.File{artifact.ArtifactMetalinkFile()},
		Generator: "boshua/boshreleasedpb",
	}

	meta4Bytes, err := metalink.Marshal(meta4)
	if err != nil {
		return errors.Wrap(err, "marshalling metalink")
	}

	return i.repository.Commit(
		map[string][]byte{path: meta4Bytes},
		fmt.Sprintf(
			"Compiling v%s for %s/%s",
			artifact.ReleaseVersion.Version,
			artifact.OSVersion.Name,
			artifact.OSVersion.Version,
		),
	)
}

func (i *index) loader() ([]compiledreleaseversion.Artifact, error) {
	paths, err := filepath.Glob(filepath.Join(i.config.RepositoryConfig.LocalPath, "compiled-release", i.config.Channel, "**", "**", "*.meta4"))
	if err != nil {
		return nil, errors.Wrap(err, "globbing")
	}

	var inmemory = []compiledreleaseversion.Artifact{}

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

		// TODO inefficient to reload; share with releaseversion.Index?
		releasePath := filepath.Join(i.config.RepositoryConfig.LocalPath, "release", i.config.Channel, path.Base(compiledReleasePath))

		releaseBytes, err := ioutil.ReadFile(releasePath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", releasePath, err)
		}

		var releaseMeta4 metalink.Metalink

		err = metalink.Unmarshal(releaseBytes, &releaseMeta4)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("unmarshalling %s: %v", releasePath, err)
		}

		// TODO sanity checks? version match? files = 1?

		inmemory = append(
			inmemory,
			compiledreleaseversion.New(
				releaseversion.Reference{
					Name:      i.config.Release,
					Version:   releaseMeta4.Files[0].Version,
					Checksums: metalinkutil.HashesToChecksums(releaseMeta4.Files[0].Hashes),
				},
				osversion.Reference{
					Name:    path.Base(path.Dir(path.Dir(compiledReleasePath))),
					Version: path.Base(path.Dir(compiledReleasePath)),
				},
				compiledReleaseMeta4.Files[0],
				map[string]interface{}{
					"uri": fmt.Sprintf(
						"%s//%s",
						i.config.RepositoryConfig.Repository,
						strings.TrimPrefix(path.Dir(strings.TrimPrefix(compiledReleasePath, i.config.RepositoryConfig.LocalPath)), "/"),
					),
					"version": compiledReleaseMeta4.Files[0].Version,
					// TODO configurable
					"options": map[string]interface{}{
						"private_key": "((index_private_key))",
					},
				},
			),
		)
	}

	i.logger.Infof("found %d entries", len(inmemory))

	return inmemory, nil
}
