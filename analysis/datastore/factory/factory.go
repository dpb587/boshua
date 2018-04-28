package factory

import (
	"fmt"

	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/analysis/datastore/presentbcr"
	analysisdatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/sirupsen/logrus"
)

type factory struct {
	logger               logrus.FieldLogger
	releaseVersionsIndex analysisdatastore.Index
}

var _ datastore.Factory = &factory{}

func New(logger logrus.FieldLogger, releaseVersionsIndex analysisdatastore.Index) datastore.Factory {
	return &factory{
		logger:               logger,
		releaseVersionsIndex: releaseVersionsIndex,
	}
}

func (f *factory) Create(provider, name string, options map[string]interface{}) (datastore.Index, error) {
	logger := f.logger.WithField("datastore", fmt.Sprintf("analysis:%s[%s]", provider, name))

	switch provider {
	case "presentbcr":
		config := presentbcr.Config{}
		err := config.Load(options)
		if err != nil {
			return nil, fmt.Errorf("loading options: %v", err)
		}

		return presentbcr.New(config, logger), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
