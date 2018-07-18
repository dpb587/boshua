package localexec

type Config struct {
	Exec string   `yaml:"exec"`
	Args []string `yaml:"args"`
}
