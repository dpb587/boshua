package boshio

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore/inmemory"
	"github.com/dpb587/metalink"

	"github.com/sirupsen/logrus"
)

const compileMatch = "bosh-stemcell-*-warden-boshlite-ubuntu-trusty-go_agent.tgz"

type index struct {
	logger   logrus.FieldLogger
	config   Config
	inmemory datastore.Index
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	idx := &index{
		logger: logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config: config,
	}

	reloader := git.NewReloader(logger, config.RepositoryConfig)

	idx.inmemory = inmemory.New(idx.loader, reloader.Reload)

	return idx
}

func (i *index) List() ([]stemcellversion.Artifact, error) {
	return i.inmemory.List()
}

func (i *index) Find(ref stemcellversion.Reference) (stemcellversion.Artifact, error) {
	return i.inmemory.Find(ref)
}

func (i *index) loader() ([]stemcellversion.Artifact, error) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/**/**/*.meta4", i.config.RepositoryConfig.LocalPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	var inmemory = []stemcellversion.Artifact{}

	for _, meta4Path := range paths {
		meta4Bytes, err := ioutil.ReadFile(meta4Path)
		if err != nil {
			return nil, fmt.Errorf("reading metalink: %v", err)
		}

		var meta4 metalink.Metalink

		err = metalink.Unmarshal(meta4Bytes, &meta4)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling metalink: %v", err)
		}

		for _, file := range meta4.Files {
			ref := ConvertFileNameToReference(file.Name)
			if ref == nil {
				// TODO log warning?
				continue
			}

			inmemory = append(
				inmemory,
				stemcellversion.New(
					*ref,
					file,
					map[string]interface{}{
						"uri": fmt.Sprintf(
							"%s//%s",
							i.config.RepositoryConfig.Repository,
							strings.TrimPrefix(path.Dir(strings.TrimPrefix(meta4Path, i.config.RepositoryConfig.LocalPath)), "/"),
						),
						"include_files": []string{
							file.Name,
						},
					},
				),
			)
		}
	}

	i.logger.Infof("found %d entries", len(inmemory))

	return inmemory, nil
}
