package boshuaV2

type BoshuaConfig struct {
	URL string `yaml:"url"`
}

func (c *BoshuaConfig) ApplyDefaults() {
	if c.URL == "" {
		c.URL = "http://127.0.0.1:4508"
	}
}
