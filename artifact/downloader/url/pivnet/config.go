package pivnet

type Config struct {
	Token      string `yaml:"token"`
	AcceptEULA bool   `yaml:"accept_eula"`
}
