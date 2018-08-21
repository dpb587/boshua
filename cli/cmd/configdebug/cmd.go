package configdebug

import (
	"fmt"

	"github.com/dpb587/boshua/config/provider/setter"
	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Cmd struct {
	setter.AppConfig `no-flag:"true"`

	Raw bool `long:"raw" description:"Show the original, raw config used"`
}

var _ flags.Commander = Cmd{}

func (c Cmd) Execute(_ []string) error {
	if c.Raw {
		if c.AppConfig.Config.Config.RawConfig == nil {
			return errors.New("raw config data is not available")
		}

		rawConfig, err := c.AppConfig.Config.Config.RawConfig()
		if err != nil {
			return errors.Wrap(err, "loading raw config")
		}

		fmt.Printf("%s\n", rawConfig)

		return nil
	}

	bytes, err := yaml.Marshal(c.AppConfig.Config.Config)
	if err != nil {
		return errors.Wrap(err, "marshaling config")
	}

	fmt.Printf("%s\n", bytes)

	return nil
}
