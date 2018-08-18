package boshioreleasesindex

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/inmemory"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type index struct {
	logger     logrus.FieldLogger
	config     Config
	repository *repository.Repository

	cache      *inmemory.Index
	cacheMutex *sync.Mutex
	cacheWarm  bool
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		logger:     logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config:     config,
		repository: repository.NewRepository(logger, config.RepositoryConfig),
		cache:      inmemory.New(),
		cacheMutex: &sync.Mutex{},
	}
}

func (i *index) GetArtifacts(f datastore.FilterParams) ([]releaseversion.Artifact, error) {
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

	paths, err := filepath.Glob(i.repository.Path("github.com", "*", "*", "*", "release.v1.yml"))
	if err != nil {
		return errors.Wrap(err, "globbing")
	}

	for _, releasePath := range paths {
		releaseBytes, err := ioutil.ReadFile(releasePath)
		if err != nil {
			return fmt.Errorf("reading %s: %v", releasePath, err)
		}

		var release releaseV1

		err = yaml.Unmarshal(releaseBytes, &release)
		if err != nil {
			// TODO warn and continue?
			return fmt.Errorf("unmarshalling %s: %v", releasePath, err)
		}

		sourcePath := filepath.Join(path.Dir(releasePath), "source.meta4")

		sourceBytes, err := ioutil.ReadFile(sourcePath)
		if err != nil {
			if os.IsNotExist(err) {
				// odd; why? e.g. github.com/cloudfoundry-incubator/diego-release/diego-0.548
				continue
			}

			return fmt.Errorf("reading %s: %v", sourcePath, err)
		}

		var sourceMeta4 metalink.Metalink

		err = metalink.Unmarshal(sourceBytes, &sourceMeta4)
		if err != nil {
			// TODO warn and continue?
			return fmt.Errorf("unmarshalling %s: %v", sourcePath, err)
		}

		sourcePathSplit := strings.Split(sourcePath, string(filepath.Separator))
		labels := append(i.config.Labels, fmt.Sprintf("repo/%s", strings.Join(sourcePathSplit[len(sourcePathSplit)-5:len(sourcePathSplit)-2], "/")))

		i.cache.Add(releaseversion.Artifact{
			Name:          release.Name,
			Version:       release.Version,
			SourceTarball: sourceMeta4.Files[0],
			Labels:        labels,
		})
	}

	i.cacheWarm = true

	return nil
}

func (i *index) GetLabels() ([]string, error) {
	err := i.fillCache()
	if err != nil {
		return nil, err
	}

	return i.cache.GetLabels()
}

func (i *index) FlushCache() error {
	i.cacheWarm = false

	// TODO defer reload?
	return i.repository.ForceReload()
}
