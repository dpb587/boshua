package boshreleasedir

import (
	"fmt"

	"github.com/dpb587/boshua/config/configdef"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const Provider = "boshreleasedir"

type factory struct {
	logger logrus.FieldLogger
}

func NewFactory(logger logrus.FieldLogger) datastore.Factory {
	return &factory{
		logger: logger,
	}
}

func (f *factory) Create(provider, name string, options map[string]interface{}) (datastore.Index, error) {
	if Provider != provider {
		return nil, fmt.Errorf("unsupported type: %s", provider)
	}

	logger := f.logger.WithField("datastore", fmt.Sprintf("releaseversion:%s[%s]", provider, name))

	cfg := Config{}
	err := configdef.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading options")
	}

	return New(cfg, logger), nil
}
