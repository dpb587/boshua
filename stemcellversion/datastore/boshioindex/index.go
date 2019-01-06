package boshioindex

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore/inmemory"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type index struct {
	name       string
	logger     logrus.FieldLogger
	config     Config
	repository *repository.Repository

	cache      *inmemory.Index
	cacheMutex *sync.Mutex
	cacheWarm  bool
}

var _ datastore.Index = &index{}

func New(name string, config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		name:       name,
		logger:     logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config:     config,
		repository: repository.NewRepository(logger, config.RepositoryConfig),
		cache:      inmemory.New(),
		cacheMutex: &sync.Mutex{},
	}
}

func (i *index) GetName() string {
	return i.name
}

func (i *index) GetArtifacts(f datastore.FilterParams, l datastore.LimitParams) ([]stemcellversion.Artifact, error) {
	err := i.fillCache()
	if err != nil {
		return nil, err
	}

	return i.cache.GetArtifacts(f, l)
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

	refs := map[stemcellversion.Reference]*stemcellversion.Artifact{}

	for t, l := range map[string]string{"published": "stability/stable", "dev": "stability/dev"} {
		paths, err := filepath.Glob(fmt.Sprintf("%s/**/**/*.meta4", i.repository.Path(t)))
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
				i.logger.Errorf("failed to unmarshal metalink: %s", meta4Path)

				continue
			}

			for _, file := range meta4.Files {
				result := ConvertFileNameToReference(file.Name)
				if result == nil {
					i.logger.Errorf("failed to extract metadata from file name: %s", file.Name)

					continue
				}

				if t == "dev" {
					// dev stemcells are created with the wrong URL :(
					for urlIdx, url := range file.URLs {
						url.URL = strings.Replace(url.URL, "https://s3.amazonaws.com/bosh-core-stemcells/", "https://s3.amazonaws.com/bosh-core-stemcells-candidate/", -1)

						file.URLs[urlIdx] = url
					}
				}

				ref := result.Reference().(stemcellversion.Reference)

				if _, found := refs[ref]; found {
					// TODO merge files and assert checksums match
					refs[ref].Labels = append(refs[ref].Labels, l)
				} else {
					result.Datastore = i.name
					result.Tarball = file
					result.Labels = append(i.config.Labels, l)

					refs[ref] = result
				}
			}
		}
	}

	for _, r := range refs {
		i.cache.Add(*r)
	}

	i.cacheWarm = true

	return nil
}

func (i *index) FlushCache() error {
	i.cacheWarm = false

	// TODO defer reload?
	return i.repository.ForceReload()
}
