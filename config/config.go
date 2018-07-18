package config

type Config struct {
	General   GeneralConfig             `yaml:"general,omitempty"`
	Scheduler *AbstractComponentConfig  `yaml:"scheduler,omitempty"`
	Stemcells []AbstractComponentConfig `yaml:"stemcell_versions"`
	Releases  []AbstractComponentConfig `yaml:"release_versions"`
	// TODO release-specific indices for compiled release datastores
	CompiledReleases []AbstractComponentConfig `yaml:"compiled_release_versions"`
	Analyses         []AbstractComponentConfig `yaml:"analyses"`
	Server           ServerConfig              `yaml:"server"`
}

type GeneralConfig struct {
	IgnoreDefaultServer bool `yaml:"ignore_default_server"`
}

type ServerConfig struct {
	Bind string          `yaml:"bind"`
	TLS  ServerTLSConfig `yaml:"tls"`
}

type ServerTLSConfig struct {
	CA          string `yaml:"ca"`
	Certificate string `yaml:"certificate"`
	PrivateKey  string `yaml:"private_key"`
}

type AbstractComponentConfig struct {
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options"`
}

func (c *Config) ApplyDefaults() {
	if c.Scheduler == nil {
		c.Scheduler = &AbstractComponentConfig{
			Type: "localexec",
			Options: map[string]interface{}{
				"exec": "boshua",
			},
		}
	}

	if c.Server.Bind == "" {
		c.Server.Bind = "127.0.0.1:4508"
	}
}
