package metalinkrepository

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"

	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/inmemory"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type index struct {
	name       string
	logger     logrus.FieldLogger
	config     Config
	repository *repository.Repository
}

var _ datastore.Index = &index{}

func New(name string, config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		name:       name,
		logger:     logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config:     config,
		repository: repository.NewRepository(logger, config.RepositoryConfig),
	}
}

func (i *index) GetName() string {
	return i.name
}

func (i *index) GetArtifacts(f datastore.FilterParams, l datastore.LimitParams) ([]releaseversion.Artifact, error) {
	if !f.NameSatisfied(i.config.Release) {
		return nil, nil
	} else if !f.LabelsSatisfied(i.config.Labels) {
		return nil, nil
	}

	err := i.repository.Reload()
	if err != nil {
		return nil, errors.Wrap(err, "reloading repository")
	}

	paths, err := filepath.Glob(i.repository.Path(i.config.Path, "*.meta4"))
	if err != nil {
		return nil, errors.Wrap(err, "globbing")
	}

	var results = []releaseversion.Artifact{}

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

		if !f.VersionSatisfied(releaseMeta4.Files[0].Version) {
			continue
		} else if !f.ChecksumSatisfied(releaseMeta4.Files[0].Hashes) {
			continue
		}

		// TODO sanity checks? version match? files = 1?
		results = append(results, releaseversion.Artifact{
			Datastore:     i.name,
			Name:          i.config.Release,
			Version:       releaseMeta4.Files[0].Version,
			SourceTarball: releaseMeta4.Files[0],
			Labels:        i.config.Labels,
		})
	}

	return inmemory.LimitArtifacts(results, l)
}

func (i *index) GetLabels() ([]string, error) {
	return i.config.Labels, nil
}

func (i *index) FlushCache() error {
	// TODO defer reload?
	return i.repository.ForceReload()
}
