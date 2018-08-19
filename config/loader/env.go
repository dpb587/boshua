package loader

import (
	"os"

	"github.com/dpb587/boshua/config"
)

// TODO
func LoadFromEnv() (*config.Config, error) {
	c := &config.Config{}

	if v := os.Getenv("BOSHUA_SERVER"); v != "" {
		c.Global.DefaultServer = v
	}

	if v := os.Getenv("BOSHUA_SERVER"); v != "" {
		c.Global.DefaultServer = v
	}

	return nil, nil
}
