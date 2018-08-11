package boshuaV2

import (
	"fmt"

	"github.com/dpb587/boshua/config/configdef"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const Provider = "boshua.v2"

type factory struct {
	logger logrus.FieldLogger
}

func NewFactory(logger logrus.FieldLogger) schedulerpkg.Factory {
	return &factory{
		logger: logger,
	}
}

func (f *factory) Create(provider string, options map[string]interface{}) (schedulerpkg.Scheduler, error) {
	if Provider != provider {
		return nil, fmt.Errorf("unsupported type: %s", provider)
	}

	cfg := Config{}
	err := configdef.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading options")
	}

	return New(cfg, f.logger), nil
}
