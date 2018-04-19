package opts

import (
	"net/http"
	"os"

	"github.com/dpb587/boshua/api/v2/client"
	"github.com/dpb587/boshua/cli/client/args"
	"github.com/sirupsen/logrus"
)

type Opts struct {
	Server      string   `long:"server" description:"API server address" default:"https://boshua.io/" env:"BOSHUA_SERVER"`
	ServerToken string   `long:"server-token" description:"API server authentication token" env:"BOSHUA_SERVER_TOKEN"`
	CACert      []string `long:"ca-cert" description:"Specific CA certificate(s) to trust" env:"BOSHUA_CA_CERT"`

	Quiet    bool          `long:"quiet" description:"Suppress informational output"`
	LogLevel args.LogLevel `long:"log-level" description:"Show additional levels of log messages" default:"FATAL" env:"BOSHUA_LOG_LEVEL"`

	logger logrus.FieldLogger
}

func (o *Opts) GetClient() *client.Client {
	return client.New(http.DefaultClient, o.Server, o.GetLogger())
}

func (o *Opts) GetLogger() logrus.FieldLogger {
	if o.logger == nil {
		panic("logger is not configured")
	}

	return o.logger
}

func (o *Opts) ConfigureLogger(command string) {
	if o.logger != nil {
		panic("logger is already configured")
	}

	var logger = logrus.New()
	logger.Out = os.Stderr
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.Level(o.LogLevel)

	o.logger = logger.WithField("cli.command", command)
}

func (o *Opts) createLogger() {

}
