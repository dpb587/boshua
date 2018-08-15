package setter

import (
	"github.com/dpb587/boshua/config/provider"
)

type AppConfig struct {
	*provider.Config
}

func (c *AppConfig) SetConfig(nc *provider.Config) {
	c.Config = nc
}
