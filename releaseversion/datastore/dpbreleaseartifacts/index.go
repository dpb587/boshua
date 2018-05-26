package dpbreleaseartifacts

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/analysis/datastore/localcache"
	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/inmemory"
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

func (i *index) Filter(ref releaseversion.Reference) ([]releaseversion.Artifact, error) {
	return i.inmemory.Filter(ref)
}

func (i *index) Find(ref releaseversion.Reference) (releaseversion.Artifact, error) {
	return i.inmemory.Find(ref)
}

func (i *index) GetAnalysisDatastore() analysisdatastore.Index {
	return localcache.New()
}

func (i *index) Store(artifact releaseversion.Artifact) error {
	artifactRef := artifact.Reference().(releaseversion.Reference)

	// TODO assert release name match?
	path := filepath.Join(
		"releaseversioniled-release",
		i.config.Channel,
		fmt.Sprintf("%s.meta4", artifactRef.Version),
	)

	meta4 := metalink.Metalink{
		Files:     []metalink.File{artifact.MetalinkFile()},
		Generator: "boshua/dpbreleaseartifacts",
	}

	meta4Bytes, err := metalink.Marshal(meta4)
	if err != nil {
		return errors.Wrap(err, "marshalling metalink")
	}

	return i.repository.Commit(
		map[string][]byte{path: meta4Bytes},
		fmt.Sprintf(
			"Add v%s",
			artifactRef.Version,
		),
	)
}

func (i *index) loader() ([]releaseversion.Artifact, error) {
	paths, err := filepath.Glob(filepath.Join(i.config.RepositoryConfig.LocalPath, "release", i.config.Channel, "*.meta4"))
	if err != nil {
		return nil, errors.Wrap(err, "globbing")
	}

	var inmemory = []releaseversion.Artifact{}

	for _, releasePath := range paths {
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

		err = metalink.Unmarshal(releaseBytes, &releaseMeta4)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("unmarshalling %s: %v", releasePath, err)
		}

		// TODO sanity checks? version match? files = 1?

		inmemory = append(
			inmemory,
			releaseversion.New(
				releaseversion.Reference{
					Name:      i.config.Release,
					Version:   releaseMeta4.Files[0].Version,
					Checksums: metalinkutil.HashesToChecksums(releaseMeta4.Files[0].Hashes),
				},
				releaseMeta4.Files[0],
			),
		)
	}

	i.logger.Infof("found %d entries", len(inmemory))

	return inmemory, nil
}
