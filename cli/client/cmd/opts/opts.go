package opts

type Opts struct {
	Server      string `long:"server" description:"Server address" default:"http://localhost:8080/" env:"CFBS_SERVER"`
	ServerToken string `long:"server-token" description:"Server authentication token" env:"CFBS_SERVER_TOKEN"`
	// CACert []string `long:"ca-cert" description:"Specific CA Certificate to trust"`
}
