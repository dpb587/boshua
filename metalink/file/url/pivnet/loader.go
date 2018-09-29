package pivnet

import (
	"fmt"
	"log"
	neturl "net/url"
	"os"
	"strconv"
	"strings"

	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/file"
	"github.com/dpb587/metalink/file/url"
	"github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logshim"
	"github.com/pkg/errors"
)

type Loader struct {
	defaultHost string
	token       string
	acceptEULA  bool
}

var _ url.Loader = &Loader{}

func NewLoader(defaultHost, token string, acceptEULA bool) url.Loader {
	return &Loader{
		defaultHost: defaultHost,
		token:       token,
		acceptEULA:  acceptEULA,
	}
}

func (f Loader) Schemes() []string {
	return []string{
		"pivnet",
	}
}

func (f Loader) Load(source metalink.URL) (file.Reference, error) {
	parsed, err := neturl.Parse(source.URL)
	if err != nil {
		return nil, errors.Wrap(err, "Parsing URI")
	}

	split := strings.Split(parsed.Path, "/")
	if len(split) != 10 {
		return nil, fmt.Errorf("expected 9 segments: %s", parsed.Path)
	} else if split[1] != "api" || split[2] != "v2" || split[3] != "products" || split[5] != "releases" || split[7] != "product_files" || split[9] != "download" {
		return nil, fmt.Errorf("invalid path: %s", parsed.Path)
	}

	productName := split[4]

	releaseId, err := strconv.Atoi(split[6])
	if err != nil {
		return nil, errors.Wrap(err, "parsing release ID")
	}

	productFileId, err := strconv.Atoi(split[8])
	if err != nil {
		return nil, errors.Wrap(err, "parsing product file ID")
	}

	hostname := parsed.Hostname()
	if hostname == "" {
		hostname = f.defaultHost
	} else {
		hostname = "https://" + hostname
	}

	client := f.newClient(hostname)

	return NewReference(client, f.acceptEULA, productName, releaseId, productFileId, parsed.Query().Get("extract")), nil
}

func (f Loader) newClient(hostname string) pivnet.Client {
	config := pivnet.ClientConfig{
		Host:      hostname,
		Token:     f.token,
		UserAgent: "boshua/0.0.0+dev", // TODO
	}

	stdoutLogger := log.New(os.Stdout, "", log.LstdFlags)
	stderrLogger := log.New(os.Stderr, "", log.LstdFlags)

	verbose := false
	logger := logshim.NewLogShim(stdoutLogger, stderrLogger, verbose)

	return pivnet.NewClient(config, logger)
}
