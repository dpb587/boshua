package concourse

import (
	"github.com/cppforlife/go-patch/patch"
)

type Config struct {
	Fly string `yaml:"fly"`

	Target   string `yaml:"target"`
	Insecure bool   `yaml:"insecure"`
	URL      string `yaml:"url"`
	Team     string `yaml:"team"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	Tasks []TaskConfig `yaml:"tasks"`
}

func (c *Config) ApplyDefaults() {
	if c.Fly == "" {
		c.Fly = "/usr/local/bin/fly"
	}
}

type TaskConfig struct {
	Type      string                 `yaml:"type"`
	Ops       []patch.OpDefinition   `yaml:"ops"`
	OpsFiles  []string               `yaml:"ops_files"`
	Vars      map[string]interface{} `yaml:"vars"`
	VarsFiles []string               `yaml:"vars_files"`
}
