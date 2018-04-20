package factory

import (
	"fmt"

	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore/legacybcr"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore/presentbcr"
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
	logger := f.logger.WithField("datastore", "dpb587/openvpn-bosh-release")

	switch provider {
	case "presentbcr":
		config := presentbcr.Config{}
		err := config.Load(options)
		if err != nil {
			return nil, fmt.Errorf("loading options: %v", err)
		}

		return presentbcr.New(config, f.releaseVersionsIndex, logger), nil
	case "legacybcr":
		config := legacybcr.Config{}
		err := config.Load(options)
		if err != nil {
			return nil, fmt.Errorf("loading options: %v", err)
		}

		return legacybcr.New(config, f.releaseVersionsIndex, logger), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
