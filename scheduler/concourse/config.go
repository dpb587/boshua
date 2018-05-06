package concourse

type ConcourseConfig struct {
	Target   string `yaml:"target"`
	Insecure bool   `yaml:"insecure"`
	URL      string `yaml:"url"`
	Team     string `yaml:"team"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	SecretsPath string `yaml:"secrets_path"`
}
