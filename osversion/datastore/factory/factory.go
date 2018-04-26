package factory

import (
	"fmt"

	"github.com/dpb587/boshua/osversion/datastore"
	"github.com/dpb587/boshua/osversion/datastore/boshio"
	"github.com/sirupsen/logrus"
)

type factory struct {
	logger logrus.FieldLogger
}

var _ datastore.Factory = &factory{}

func New(logger logrus.FieldLogger) datastore.Factory {
	return &factory{
		logger: logger,
	}
}

func (f *factory) Create(provider, name string, options map[string]interface{}) (datastore.Index, error) {
	logger := f.logger.WithField("datastore", name)

	switch provider {
	case "boshio":
		config := boshio.Config{}
		err := config.Load(options)
		if err != nil {
			return nil, fmt.Errorf("loading options: %v", err)
		}

		return boshio.New(config, logger), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
