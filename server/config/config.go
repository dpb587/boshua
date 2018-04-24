package config

type Config struct {
	Concourse        ConcourseConfig   `yaml:"concourse"`
	Stemcells        []DatastoreConfig `yaml:"stemcell_versions"`
	Releases         []DatastoreConfig `yaml:"release_versions"`
	CompiledReleases []DatastoreConfig `yaml:"compiled_release_versions"`
}

type DatastoreConfig struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Options map[string]interface{}
}

type ConcourseConfig struct {
	Target   string `yaml:"target"`
	Insecure bool   `yaml:"insecure"`
	URL      string `yaml:"url"`
	Team     string `yaml:"team"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	SecretsPath string `yaml:"secrets_path"`
}
