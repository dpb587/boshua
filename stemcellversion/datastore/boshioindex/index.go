package boshioindex

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/git"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore/inmemory"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger     logrus.FieldLogger
	config     Config
	repository *git.Repository

	cache      *inmemory.Index
	cacheMutex *sync.Mutex
	cacheWarm  bool
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		logger:     logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config:     config,
		repository: git.NewRepository(logger, config.RepositoryConfig),
		cache:      inmemory.New(),
		cacheMutex: &sync.Mutex{},
	}
}

func (i *index) GetArtifacts(f datastore.FilterParams) ([]stemcellversion.Artifact, error) {
	err := i.fillCache()
	if err != nil {
		return nil, err
	}

	return i.cache.GetArtifacts(f)
}

func (i *index) fillCache() error {
	i.cacheMutex.Lock()
	defer i.cacheMutex.Unlock()

	if i.cacheWarm && i.repository.WarmCache() {
		return nil
	}

	err := i.cache.FlushCache()
	if err != nil {
		return errors.Wrap(err, "flushing in-memory")
	}

	err = i.repository.Reload()
	if err != nil {
		return errors.Wrap(err, "reloading repository")
	}

	paths, err := filepath.Glob(fmt.Sprintf("%s/**/**/*.meta4", i.repository.Path(i.config.Prefix)))
	if err != nil {
		return errors.Wrap(err, "globbing")
	}

	for _, meta4Path := range paths {
		meta4Bytes, err := ioutil.ReadFile(meta4Path)
		if err != nil {
			return errors.Wrap(err, "reading metalink")
		}

		var meta4 metalink.Metalink

		err = metalink.Unmarshal(meta4Bytes, &meta4)
		if err != nil {
			return errors.Wrap(err, "unmarshaling metalink")
		}

		for _, file := range meta4.Files {
			result := ConvertFileNameToReference(file.Name)
			if result == nil {
				// TODO log warning?
				continue
			}

			result.Tarball = file
			result.Labels = i.config.Labels

			i.cache.Add(*result)
		}
	}

	i.cacheWarm = true

	return nil
}

func (i *index) FlushCache() error {
	i.cacheWarm = false

	// TODO defer reload?
	return i.repository.ForceReload()
}
