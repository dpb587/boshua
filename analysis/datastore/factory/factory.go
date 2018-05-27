package factory

import (
	"fmt"

	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/analysis/datastore/dpbreleaseartifacts"
	"github.com/dpb587/boshua/analysis/datastore/localcache"
	"github.com/dpb587/boshua/config"
	"github.com/pkg/errors"
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
	logger := f.logger.WithField("datastore", fmt.Sprintf("analysis:%s[%s]", provider, name))

	switch provider {
	case "dpbreleaseartifacts":
		cfg := boshreleasedpb.Config{}
		err := config.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, errors.Wrap(err, "loading options")
		}

		return boshreleasedpb.New(cfg, logger), nil
	case "localcache":
		return localcache.New(), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
