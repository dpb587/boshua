package factory

import (
	"fmt"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore/boshio"
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

func (f *factory) Create(provider, name string, options map[string]interface{}, analysisIndex analysisdatastore.Index) (datastore.Index, error) {
	logger := f.logger.WithField("datastore", fmt.Sprintf("stemcellversion:%s[%s]", provider, name))

	switch provider {
	case "boshio":
		config := boshio.Config{}
		err := config.Load(options)
		if err != nil {
			return nil, errors.Wrap(err, "loading options")
		}

		return boshio.New(config, analysisIndex, logger), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
