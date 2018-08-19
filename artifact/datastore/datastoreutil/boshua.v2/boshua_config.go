package boshuaV2

// BoshuaConfig defines a remote server for boshua.v2 API.
type BoshuaConfig struct {
	// URL must define the remote HTTP endpoint.
	URL string `yaml:"url"`
}

func (c *BoshuaConfig) ApplyDefaults() {
	if c.URL == "" {
		c.URL = "http://127.0.0.1:4508"
	}
}
