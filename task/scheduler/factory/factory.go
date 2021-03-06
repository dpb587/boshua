package factory

import (
	"fmt"

	"github.com/dpb587/boshua/config/provider"
	"github.com/dpb587/boshua/task/scheduler"
	boshuaV2 "github.com/dpb587/boshua/task/scheduler/boshua.v2"
	"github.com/dpb587/boshua/task/scheduler/concourse"
	"github.com/dpb587/boshua/task/scheduler/localexec"
	"github.com/dpb587/boshua/util/configdef"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type factory struct {
	config *provider.Config
	logger logrus.FieldLogger
}

var _ scheduler.Factory = &factory{}

func New(config *provider.Config, logger logrus.FieldLogger) scheduler.Factory {
	return &factory{
		config: config,
		logger: logger,
	}
}

func (f *factory) Create(provider string, options map[string]interface{}) (scheduler.Scheduler, error) {
	logger := f.logger.WithField("scheduler", provider)

	switch provider {
	case "concourse":
		cfg := concourse.Config{}
		err := configdef.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, errors.Wrap(err, "loading options")
		}

		return concourse.New(cfg, f.config, logger), nil
	case "boshua.v2":
		cfg := boshuaV2.Config{}
		err := configdef.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, errors.Wrap(err, "loading options")
		}

		return boshuaV2.New(cfg, logger), nil
	case "localexec":
		cfg := localexec.Config{}
		err := configdef.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, errors.Wrap(err, "loading options")
		}

		return localexec.New(cfg, f.config, logger), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
