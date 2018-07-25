package boshioreleasesindex

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type Index struct {
	logger     logrus.FieldLogger
	config     Config
	repository *git.Repository
}

var _ datastore.Index = &Index{}

func New(config Config, logger logrus.FieldLogger) *Index {
	return &Index{
		logger:     logger.WithField("build.package", reflect.TypeOf(Index{}).PkgPath()),
		config:     config,
		repository: git.NewRepository(logger, config.RepositoryConfig),
	}
}

func (i *Index) Filter(f *datastore.FilterParams) ([]releaseversion.Artifact, error) {
	if !f.LabelsSatisfied(i.config.Labels) {
		return nil, nil
	}

	err := i.repository.Reload()
	if err != nil {
		return nil, errors.Wrap(err, "reloading repository")
	}

	paths, err := filepath.Glob(i.repository.Path("github.com", "*", "*", "*", "release.v1.yml"))
	if err != nil {
		return nil, errors.Wrap(err, "globbing")
	}

	var results = []releaseversion.Artifact{}

	for _, releasePath := range paths {
		releaseBytes, err := ioutil.ReadFile(releasePath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", releasePath, err)
		}

		var release releaseV1

		err = yaml.Unmarshal(releaseBytes, &release)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("unmarshalling %s: %v", releasePath, err)
		}

		if !f.NameSatisfied(release.Name) {
			continue
		} else if !f.VersionSatisfied(release.Version) {
			continue
		}

		sourcePath := filepath.Join(path.Dir(releasePath), "source.meta4")

		sourceBytes, err := ioutil.ReadFile(sourcePath)
		if err != nil {
			if os.IsNotExist(err) {
				// odd; why? e.g. github.com/cloudfoundry-incubator/diego-release/diego-0.548
				continue
			}

			return nil, fmt.Errorf("reading %s: %v", sourcePath, err)
		}

		var sourceMeta4 metalink.Metalink

		err = metalink.Unmarshal(sourceBytes, &sourceMeta4)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("unmarshalling %s: %v", sourcePath, err)
		}

		if !f.ChecksumSatisfied(sourceMeta4.Files[0].Hashes) {
			continue
		}

		sourcePathSplit := strings.Split(sourcePath, string(filepath.Separator))

		// TODO sanity checks? version match? files = 1?
		results = append(results, releaseversion.Artifact{
			Name:          release.Name,
			Version:       release.Version,
			SourceTarball: sourceMeta4.Files[0],
			Labels:        append(i.config.Labels, fmt.Sprintf("repo/%s", strings.Join(sourcePathSplit[len(sourcePathSplit)-5:len(sourcePathSplit)-2], "/"))),
		})
	}

	return results, nil
}
