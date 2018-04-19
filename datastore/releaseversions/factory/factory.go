package factory

import (
	"fmt"

	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions/boshio"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions/meta4"
	"github.com/sirupsen/logrus"
)

type factory struct {
	logger logrus.FieldLogger
}

var _ releaseversions.Factory = &factory{}

func New(logger logrus.FieldLogger) releaseversions.Factory {
	return &factory{
		logger: logger,
	}
}

func (f *factory) Create(provider, name string, options map[string]interface{}) (releaseversions.Index, error) {
	logger := f.logger.WithField("datastore", name)

	switch provider {
	case "boshio":
		config := boshio.Config{}
		err := config.Load(options)
		if err != nil {
			return nil, fmt.Errorf("loading options: %v", err)
		}

		return boshio.New(config, logger), nil
	case "meta4":
		config := meta4.Config{}
		err := config.Load(options)
		if err != nil {
			return nil, fmt.Errorf("loading options: %v", err)
		}

		return meta4.New(config, logger), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
