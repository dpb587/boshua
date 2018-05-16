package boshreleasedir

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/analysis/datastore/localcache"
	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/metalink/file/metaurl/boshreleasesource"
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

func (i *index) loader() ([]releaseversion.Artifact, error) {
	indices, err := filepath.Glob(fmt.Sprintf("%s/releases/**/index.yml", i.config.RepositoryConfig.LocalPath))
	if err != nil {
		return nil, errors.Wrap(err, "globbing")
	}

	var inmemory = []releaseversion.Artifact{}

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
			var release boshRelease

			releaseName := path.Base(path.Dir(indexPath))
			releaseSubPath := fmt.Sprintf("releases/%s/%s-%s.yml", releaseName, releaseName, build.Version)
			releasePath := filepath.Join(i.config.RepositoryConfig.LocalPath, releaseSubPath)

			releaseBytes, err := ioutil.ReadFile(releasePath)
			if err != nil {
				return nil, fmt.Errorf("reading %s: %v", releasePath, err)
			}

			releaseSha1 := sha1.New()
			releaseSha1.Write(releaseBytes)
			sha1 := fmt.Sprintf("%x", releaseSha1.Sum(nil))

			err = yaml.Unmarshal(releaseBytes, &release)
			if err != nil {
				return nil, fmt.Errorf("unmarshalling %s: %v", indexPath, err)
			}

			ref := releaseversion.Reference{
				Name:    release.Name,
				Version: release.Version,
				URLs: []string{
					fmt.Sprintf("%s//%s#%s", i.config.RepositoryConfig.Repository, releaseSubPath, sha1),
				},
			}

			inmemory = append(inmemory, releaseversion.New(
				ref,
				metalink.File{
					Name: fmt.Sprintf("%s-%s.tgz", ref.Name, ref.Version),
					MetaURLs: []metalink.MetaURL{
						{
							URL:       fmt.Sprintf("%s//%s", i.config.RepositoryConfig.Repository, releaseSubPath),
							MediaType: boshreleasesource.DefaultMediaType,
						},
					},
				},
			))
		}
	}

	i.logger.Infof("found %d entries", len(inmemory))

	return inmemory, nil
}
