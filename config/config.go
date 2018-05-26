package config

type Config struct {
	General          GeneralConfig             `yaml:"general,omitempty"`
	Scheduler        *AbstractComponentConfig  `yaml:"scheduler,omitempty"`
	Stemcells        []AbstractComponentConfig `yaml:"stemcell_versions"`
	Releases         []AbstractComponentConfig `yaml:"release_versions"`
	CompiledReleases []AbstractComponentConfig `yaml:"compiled_release_versions"`
	Analyses         []AbstractComponentConfig `yaml:"analyses"`
	Server           ServerConfig              `yaml:"server"`
}

type GeneralConfig struct {
	IgnoreDefaultServer bool `yaml:"ignore_default_server"`
}

// TODO
type ServerConfig struct {
	Certificate string `yaml:"certificate"`
	PrivateKey  string `yaml:"private_key"`
}

type AbstractComponentConfig struct {
	Name     string                   `yaml:"name"`
	Type     string                   `yaml:"type"`
	Options  map[string]interface{}   `yaml:"options"`
	Analysis *AbstractComponentConfig `yaml:"analysis"`
}

func (c *Config) ApplyDefaults() {
	if c.Scheduler == nil {
		c.Scheduler = &AbstractComponentConfig{
			Type: "localexec",
		}
	}
}
