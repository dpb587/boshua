package opts

import (
	"github.com/dpb587/boshua/cli/opts"
	"github.com/dpb587/boshua/config/presets/defaults"
	"github.com/dpb587/boshua/config/provider"
)

type Opts struct {
	*opts.Opts

	parsedConfig *provider.Config
}

func (o *Opts) GetConfig() (*provider.Config, error) {
	if o.parsedConfig == nil {
		cfg, err := o.Opts.GetConfig()
		if err != nil {
			return nil, err
		}

		defaults.UseDefaultFactories(cfg)

		o.parsedConfig = cfg
	}

	return o.parsedConfig, nil
}
