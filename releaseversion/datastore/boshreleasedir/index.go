package boshreleasedir

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"

	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/metalink/file/metaurl/boshreleasesource"
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
	if f.ChecksumExpected {
		return nil, nil
	}

	var globReleaseName = "*"
	if f.NameExpected {
		globReleaseName = filepath.Base(f.Name)
	}

	indices, err := filepath.Glob(filepath.Join(i.config.RepositoryConfig.LocalPath, "releases", globReleaseName, "index.yml"))
	if err != nil {
		return nil, errors.Wrap(err, "globbing")
	}

	var results = []releaseversion.Artifact{}

	for _, indexPath := range indices {
		indexBytes, err := ioutil.ReadFile(indexPath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", indexPath, err)
		}

		var index boshReleaseIndex

		err = yaml.Unmarshal(indexBytes, &index)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling %s: %v", indexPath, err)
		}

		for _, build := range index.Builds {
			releaseName := path.Base(path.Dir(indexPath))
			releaseSubPath := filepath.Join("releases", releaseName, fmt.Sprintf("%s-%s.yml", releaseName, build.Version))
			releasePath := filepath.Join(i.config.RepositoryConfig.LocalPath, releaseSubPath)

			// TODO track with checksum of release manifest somewhere?
			releaseBytes, err := ioutil.ReadFile(releasePath)
			if err != nil {
				return nil, fmt.Errorf("reading %s: %v", releasePath, err)
			}

			var release boshRelease

			err = yaml.Unmarshal(releaseBytes, &release)
			if err != nil {
				return nil, fmt.Errorf("unmarshalling %s: %v", indexPath, err)
			}

			if !f.NameSatisfied(release.Name) {
				continue
			} else if !f.VersionSatisfied(release.Version) {
				continue
			}

			metaurls := []metalink.MetaURL{
				{
					URL:       fmt.Sprintf("%s//%s", i.config.RepositoryConfig.Repository, releaseSubPath),
					MediaType: boshreleasesource.DefaultMediaType,
				},
			}

			if !f.URISatisfied(nil, metaurls) {
				continue
			}

			results = append(results, releaseversion.Artifact{
				Name:    release.Name,
				Version: release.Version,
				SourceTarball: metalink.File{
					Name:     fmt.Sprintf("%s-%s.tgz", release.Name, release.Version),
					MetaURLs: metaurls,
				},
			})
		}
	}

	return results, nil
}
