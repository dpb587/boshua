package cfdeployment

import (
	"fmt"

	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/util/configdef"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const ProviderName datastore.ProviderName = "cfdeployment"

type factory struct {
	releaseVersionIndex releaseversiondatastore.NamedGetter
	logger              logrus.FieldLogger
}

func NewFactory(releaseVersionIndex releaseversiondatastore.NamedGetter, logger logrus.FieldLogger) datastore.Factory {
	return &factory{
		releaseVersionIndex: releaseVersionIndex,
		logger:              logger,
	}
}

func (f *factory) Create(provider datastore.ProviderName, name string, options map[string]interface{}) (datastore.Index, error) {
	if ProviderName != provider {
		return nil, fmt.Errorf("unsupported type: %s", provider)
	}

	logger := f.logger.WithField("datastore", fmt.Sprintf("compilation:%s[%s]", provider, name))

	cfg := Config{}
	err := configdef.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading options")
	}

	releaseVersionIndex, err := f.releaseVersionIndex("default") // TODO configurable
	if err != nil {
		return nil, errors.Wrap(err, "loading release index")
	}

	return New(name, releaseVersionIndex, cfg, logger), nil
}
