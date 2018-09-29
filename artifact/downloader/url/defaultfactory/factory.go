package factory

import (
	"github.com/dpb587/boshua/artifact/downloader/url"
	"github.com/dpb587/boshua/artifact/downloader/url/s3"
	"github.com/dpb587/boshua/artifact/downloader/url/pivnet"
	"github.com/dpb587/boshua/artifact/downloader/url/factory"
	"github.com/sirupsen/logrus"
)

func New(logger logrus.FieldLogger) url.Factory {
	f := factory.New()
	f.Add(pivnet.ProviderName, pivnet.NewFactory(logger))
	f.Add(s3.ProviderName, s3.NewFactory(logger))

	return f
}
