package factory

import (
	"fmt"

	"github.com/dpb587/boshua/config"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/dpb587/boshua/task/scheduler/concourse"
	"github.com/dpb587/boshua/task/scheduler/localexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type factory struct {
	logger logrus.FieldLogger
}

var _ scheduler.Factory = &factory{}

func New(logger logrus.FieldLogger) scheduler.Factory {
	return &factory{
		logger: logger,
	}
}

func (f *factory) Create(provider string, options map[string]interface{}) (scheduler.Scheduler, error) {
	logger := f.logger.WithField("scheduler", provider)

	switch provider {
	case "concourse":
		cfg := concourse.Config{}
		err := config.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, errors.Wrap(err, "loading options")
		}

		return concourse.New(cfg, logger), nil
	case "localexec":
		cfg := localexec.Config{}
		err := config.RemarshalYAML(options, &cfg)
		if err != nil {
			return nil, errors.Wrap(err, "loading options")
		}

		return localexec.New(cfg, logger), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
