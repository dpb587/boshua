package factory

import (
	"fmt"

	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore/contextualosmetalinkrepository"
	// "github.com/dpb587/boshua/releaseversion/compilation/datastore/legacybcr"
	"github.com/dpb587/boshua/config"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type factory struct {
	logger               logrus.FieldLogger
	releaseVersionsIndex releaseversiondatastore.Index
}

var _ datastore.Factory = &factory{}

func New(logger logrus.FieldLogger) datastore.Factory {
	return &factory{
		logger: logger,
	}
}

func (f *factory) Create(provider, name string, options map[string]interface{}, releaseVersionIndex releaseversiondatastore.Index) (datastore.Index, error) {
	logger := f.logger.WithField("datastore", fmt.Sprintf("compiledreleaseversion:%s[%s]", provider, name))

	switch provider {
	case "contextualosmetalinkrepository":
		cfg := contextualosmetalinkrepository.Config{}
		err := config.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, errors.Wrap(err, "loading options")
		}

		return contextualosmetalinkrepository.New(releaseVersionIndex, cfg, logger), nil
	// case "presentbcr":
	// 	cfg := presentbcr.Config{}
	// 	err := config.RemarshalYAML(options, &cfg)
	// 	if err != nil {
	// 		return nil, errors.Wrap(err, "loading options")
	// 	}
	//
	// 	return presentbcr.New(cfg, logger), nil
	// case "legacybcr":
	// 	cfg := legacybcr.Config{}
	// 	err := config.RemarshalYAML(options, &cfg)
	// 	if err != nil {
	// 		return nil, errors.Wrap(err, "loading options")
	// 	}
	//
	// 	return legacybcr.New(cfg, f.releaseVersionsIndex, logger), nil
	// case "compiledreleasesops": // TODO https://github.com/cloudfoundry/cf-deployment/blob/master/operations/use-compiled-releases.yml
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
