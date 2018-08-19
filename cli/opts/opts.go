package opts

import (
	"time"

	"github.com/dpb587/boshua/cli/args"
	"github.com/dpb587/boshua/config"
	configloader "github.com/dpb587/boshua/config/loader"
	configprovider "github.com/dpb587/boshua/config/provider"
	configtypes "github.com/dpb587/boshua/config/types"
	"github.com/pkg/errors"
)

type Opts struct {
	Config string `long:"config" description:"Path to configuration file" env:"BOSHUA_CONFIG" default:"~/.config/boshua/config.yml"`

	DefaultServer string        `long:"default-server" description:"Default boshua API server" env:"BOSHUA_SERVER"`
	DefaultWait   args.Duration `long:"default-wait" description:"Maximum time to wait for scheduled tasks; 0 to disable scheduling" env:"BOSHUA_WAIT" default:"30m"` // TODO better name

	Quiet    bool          `long:"quiet" description:"Suppress informational output" env:"BOSHUA_QUIET"`
	LogLevel args.LogLevel `long:"log-level" description:"Show additional levels of log messages" default:"FATAL" env:"BOSHUA_LOG_LEVEL"`

	parsedConfig *configprovider.Config
}

func (o *Opts) GetConfig() (*configprovider.Config, error) {
	if o.parsedConfig != nil {
		return o.parsedConfig, nil
	}

	cfg, err := configloader.LoadFromFile(
		o.Config,
		&config.Config{
			Global: config.GlobalConfig{
				DefaultServer: o.DefaultServer,
				DefaultWait:   time.Duration(o.DefaultWait),
				LogLevel:      configtypes.LogLevel(o.LogLevel),
				Quiet:         o.Quiet,
			},
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "loading file")
	}

	o.parsedConfig = cfg

	return o.parsedConfig, nil
}
