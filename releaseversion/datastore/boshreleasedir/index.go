package boshreleasedir

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"

	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/inmemory"
	"github.com/dpb587/metalink"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type index struct {
	logger     logrus.FieldLogger
	repository string
	localPath  string
	inmemory   datastore.Index
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	idx := &index{
		logger:     logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		repository: config.Repository,
		localPath:  config.LocalPath,
	}

	reloader := git.NewReloader(logger, config.Repository, config.LocalPath, config.PullInterval)

	idx.inmemory = inmemory.New(idx.loader, reloader.Reload)

	return idx
}

func (i *index) List() ([]releaseversion.Artifact, error) {
	return i.inmemory.List()
}

func (i *index) Find(ref releaseversion.Reference) (releaseversion.Artifact, error) {
	return i.inmemory.Find(ref)
}

func (i *index) loader() ([]releaseversion.Artifact, error) {
	indices, err := filepath.Glob(fmt.Sprintf("%s/releases/**/index.yml", i.localPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
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
			releasePath := filepath.Join(i.localPath, releaseSubPath)

			releaseBytes, err := ioutil.ReadFile(releasePath)
			if err != nil {
				return nil, fmt.Errorf("reading %s: %v", releasePath, err)
			}

			err = yaml.Unmarshal(releaseBytes, &release)
			if err != nil {
				return nil, fmt.Errorf("unmarshalling %s: %v", indexPath, err)
			}

			ref := releaseversion.Reference{
				Name:    release.Name,
				Version: release.Version,
			}

			inmemory = append(inmemory, releaseversion.New(
				ref,
				metalink.File{
					Name: fmt.Sprintf("%s-%s.tgz", ref.Name, ref.Version),
					MetaURLs: []metalink.MetaURL{
						{
							URL:       fmt.Sprintf("%s//%s", i.repository, releaseSubPath),
							MediaType: "application/vnd.bosh.release",
						},
					},
				},
				map[string]interface{}{},
			))
		}
	}

	i.logger.Infof("found %d entries", len(inmemory))

	return inmemory, nil
}
