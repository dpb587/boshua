package concourse

type Config struct {
	Fly string `yaml:"fly"`

	Target   string `yaml:"target"`
	Insecure bool   `yaml:"insecure"`
	URL      string `yaml:"url"`
	Team     string `yaml:"team"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	SecretsPath string `yaml:"secrets_path"` // TODO remove?
}

func (c *Config) ApplyDefaults() {
	if c.Fly == "" {
		c.Fly = "/usr/local/bin/fly"
	}
}
