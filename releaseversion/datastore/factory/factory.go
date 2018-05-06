package factory

import (
	"fmt"

	"github.com/dpb587/boshua/config"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/boshio"
	"github.com/dpb587/boshua/releaseversion/datastore/boshreleasedir"
	"github.com/dpb587/boshua/releaseversion/datastore/meta4"
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
	logger := f.logger.WithField("datastore", fmt.Sprintf("releaseversion:%s[%s]", provider, name))

	switch provider {
	case "boshio":
		cfg := boshio.Config{}
		err := config.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, fmt.Errorf("loading options: %v", err)
		}

		return boshio.New(cfg, logger), nil
	case "boshreleasedir":
		cfg := boshreleasedir.Config{}
		err := config.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, fmt.Errorf("loading options: %v", err)
		}

		return boshreleasedir.New(cfg, logger), nil
	case "meta4":
		cfg := meta4.Config{}
		err := config.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, fmt.Errorf("loading options: %v", err)
		}

		return meta4.New(cfg, logger), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
