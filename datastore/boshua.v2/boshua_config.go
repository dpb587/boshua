package boshuav2

type BoshuaConfig struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	CACert string `yaml:"ca_cert"`
}
