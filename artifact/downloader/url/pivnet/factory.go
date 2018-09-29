package pivnet

import (
	"fmt"

	downloaderurl "github.com/dpb587/boshua/artifact/downloader/url"
	"github.com/dpb587/metalink/file/url"
	"github.com/dpb587/boshua/metalink/file/url/pivnet"
	"github.com/dpb587/boshua/util/configdef"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const ProviderName = "pivnet"

type factory struct {
	logger logrus.FieldLogger
}

func NewFactory(logger logrus.FieldLogger) downloaderurl.Factory {
	return &factory{
		logger: logger,
	}
}

func (f *factory) Create(provider downloaderurl.ProviderName, name string, options map[string]interface{}) (url.Loader, error) {
	if ProviderName != provider {
		return nil, fmt.Errorf("unsupported type: %s", provider)
	}

	cfg := Config{}
	err := configdef.RemarshalYAML(options, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "loading options")
	}

	return pivnet.NewLoader(cfg.Token, cfg.AcceptEULA), nil
}
