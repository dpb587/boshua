package boshioindex

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"

	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger     logrus.FieldLogger
	config     Config
	repository *git.Repository
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		logger:     logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config:     config,
		repository: git.NewRepository(logger, config.RepositoryConfig),
	}
}

func (i *index) Filter(f *datastore.FilterParams) ([]stemcellversion.Artifact, error) {
	if !f.LabelsSatisfied(i.config.Labels) {
		return nil, nil
	}

	err := i.repository.Reload()
	if err != nil {
		return nil, errors.Wrap(err, "reloading repository")
	}

	paths, err := filepath.Glob(fmt.Sprintf("%s/**/**/*.meta4", i.repository.Path(i.config.Prefix)))
	if err != nil {
		return nil, errors.Wrap(err, "globbing")
	}

	var results = []stemcellversion.Artifact{}

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
			result := ConvertFileNameToReference(file.Name)
			if result == nil {
				// TODO log warning?
				continue
			}

			if !f.OSSatisfied(result.OS) {
				continue
			} else if !f.VersionSatisfied(result.Version) {
				continue
			} else if !f.IaaSSatisfied(result.IaaS) {
				continue
			} else if !f.HypervisorSatisfied(result.Hypervisor) {
				continue
			} else if !f.FlavorSatisfied(result.Flavor) {
				continue
			}

			result.Tarball = file
			result.Labels = i.config.Labels

			results = append(results, *result)
		}
	}

	return results, nil
}
