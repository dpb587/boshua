package factory

import (
	"fmt"

	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore/boshreleasedpb"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore/legacybcr"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore/presentbcr"
	"github.com/dpb587/boshua/config"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/sirupsen/logrus"
)

type factory struct {
	logger               logrus.FieldLogger
	releaseVersionsIndex releaseversiondatastore.Index
}

var _ datastore.Factory = &factory{}

func New(logger logrus.FieldLogger, releaseVersionsIndex releaseversiondatastore.Index) datastore.Factory {
	return &factory{
		logger:               logger,
		releaseVersionsIndex: releaseVersionsIndex,
	}
}

func (f *factory) Create(provider, name string, options map[string]interface{}) (datastore.Index, error) {
	logger := f.logger.WithField("datastore", fmt.Sprintf("compiledreleaseversion:%s[%s]", provider, name))

	switch provider {
	case "boshreleasedpb":
		cfg := boshreleasedpb.Config{}
		err := config.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, fmt.Errorf("loading options: %v", err)
		}

		return boshreleasedpb.New(cfg, logger), nil
	case "presentbcr":
		cfg := presentbcr.Config{}
		err := config.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, fmt.Errorf("loading options: %v", err)
		}

		return presentbcr.New(cfg, logger), nil
	case "legacybcr":
		cfg := legacybcr.Config{}
		err := config.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, fmt.Errorf("loading options: %v", err)
		}

		return legacybcr.New(cfg, f.releaseVersionsIndex, logger), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
