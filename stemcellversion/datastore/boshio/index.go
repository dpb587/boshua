package boshio

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore/inmemory"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const compileMatch = "bosh-stemcell-*-warden-boshlite-ubuntu-trusty-go_agent.tgz"

type index struct {
	logger        logrus.FieldLogger
	config        Config
	repository    *git.Repository
	inmemory      datastore.Index
	analysisIndex analysisdatastore.Index
}

var _ datastore.Index = &index{}

func New(config Config, analysisIndex analysisdatastore.Index, logger logrus.FieldLogger) datastore.Index {
	idx := &index{
		logger:        logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config:        config,
		repository:    git.NewRepository(logger, config.RepositoryConfig),
		analysisIndex: analysisIndex,
	}

	idx.inmemory = inmemory.New(idx.loader, idx.repository.Reload)

	return idx
}

func (i *index) Find(ref stemcellversion.Reference) (stemcellversion.Artifact, error) {
	return i.inmemory.Find(ref)
}

func (i *index) Filter(ref stemcellversion.Reference) ([]stemcellversion.Artifact, error) {
	return i.inmemory.Filter(ref)
}

func (i *index) List() ([]stemcellversion.Artifact, error) {
	return i.inmemory.List()
}

func (i *index) GetAnalysisDatastore(ref stemcellversion.Reference) (analysisdatastore.Index, error) {
	if i.analysisIndex == nil {
		return nil, datastore.UnsupportedOperationErr
	}

	return i.analysisIndex, nil
}

func (i *index) loader() ([]stemcellversion.Artifact, error) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/**/**/*.meta4", i.config.RepositoryConfig.LocalPath))
	if err != nil {
		return nil, errors.Wrap(err, "globbing")
	}

	var inmemory = []stemcellversion.Artifact{}

	for _, meta4Path := range paths {
		meta4Bytes, err := ioutil.ReadFile(meta4Path)
		if err != nil {
			return nil, errors.Wrap(err, "reading metalink")
		}

		var meta4 metalink.Metalink

		err = metalink.Unmarshal(meta4Bytes, &meta4)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshaling metalink")
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
				),
			)
		}
	}

	i.logger.Infof("found %d entries", len(inmemory))

	return inmemory, nil
}
