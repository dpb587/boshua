package pivnet

import (
	"github.com/pivotal-cf/go-pivnet"
)

type Config struct {
	Host  string `yaml:"host"`
	Token string `yaml:"token"`
}

func (c *Config) ApplyDefaults() {
	if c.Host == "" {
		c.Host = pivnet.DefaultHost
	}
}
