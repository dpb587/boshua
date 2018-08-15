package loader

import (
	"os"

	"github.com/dpb587/boshua/config"
)

// TODO
func LoadFromEnv() (*config.Config, error) {
	c := &config.Config{}

	if v := os.Getenv("BOSHUA_SERVER"); v != "" {
		c.General.DefaultServer = v
	}

	if v := os.Getenv("BOSHUA_SERVER"); v != "" {
		c.General.DefaultServer = v
	}

	return nil, nil
}
