package boshioindex

import (
	"fmt"

	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/util/configdef"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const ProviderName = "boshioindex"

type factory struct {
	logger logrus.FieldLogger
}

func NewFactory(logger logrus.FieldLogger) datastore.Factory {
	return &factory{
		logger: logger,
	}
}

func (f *factory) Create(provider datastore.ProviderName, name string, options map[string]interface{}) (datastore.Index, error) {
	if ProviderName != provider {
		return nil, fmt.Errorf("unsupported type: %s", provider)
	}

	logger := f.logger.WithField("datastore", fmt.Sprintf("stemcellversion:%s[%s]", provider, name))

	cfg := Config{}
	err := configdef.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading options")
	}

	return New(cfg, logger), nil
}
