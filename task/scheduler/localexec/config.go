package localexec

type Config struct {
	Exec string   `yaml:"exec"`
	Args []string `yaml:"args"`
}

func (c *Config) ApplyDefaults() {
	if c.Exec == "" {
		c.Exec = "boshua"
	}
}
