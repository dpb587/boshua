package factory

import (
	"fmt"

	"github.com/dpb587/boshua/config/configdef"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/boshioreleasesindex"
	"github.com/dpb587/boshua/releaseversion/datastore/boshreleasedir"
	boshuaV2 "github.com/dpb587/boshua/releaseversion/datastore/boshua.v2"
	"github.com/dpb587/boshua/releaseversion/datastore/metalinkrepository"
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
	logger := f.logger.WithField("datastore", fmt.Sprintf("releaseversion:%s[%s]", provider, name))

	switch provider {
	case "boshioreleasesindex":
		cfg := boshioreleasesindex.Config{}
		err := configdef.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, errors.Wrap(err, "loading options")
		}

		return boshioreleasesindex.New(cfg, logger), nil
	case "boshua.v2":
		cfg := boshuaV2.Config{}
		err := configdef.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, errors.Wrap(err, "loading options")
		}

		return boshuaV2.New(cfg, logger), nil
	case "boshreleasedir":
		cfg := boshreleasedir.Config{}
		err := configdef.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, errors.Wrap(err, "loading options")
		}

		return boshreleasedir.New(cfg, logger), nil
	case "metalinkrepository":
		cfg := metalinkrepository.Config{}
		err := configdef.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, errors.Wrap(err, "loading options")
		}

		return metalinkrepository.New(cfg, logger), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
