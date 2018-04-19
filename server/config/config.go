package config

type Config struct {
	Concourse        ConcourseConfig   `yaml:"concourse"`
	Stemcells        []DatastoreConfig `yaml:"stemcells"`
	Releases         []DatastoreConfig `yaml:"releases"`
	CompiledReleases []DatastoreConfig `yaml:"compiled_releases"`
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

	PipelinePath string `yaml:"pipeline_path"`
	SecretsPath  string `yaml:"vars_path"`
}
