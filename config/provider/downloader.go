package provider

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/boshua/metalink/file/metaurl/boshreleasesource"
	"github.com/dpb587/metalink/file/url"
	"github.com/dpb587/metalink/file/metaurl"
	urldefaultloader "github.com/dpb587/metalink/file/url/defaultloader"
)

func (c *Config) GetDownloader() (url.Loader, metaurl.Loader, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)

	urlLoader := urldefaultloader.New(fs)
	metaurlLoader := metaurl.NewLoaderFactory()
	metaurlLoader.Add(boshreleasesource.Loader{})

	return urlLoader, metaurlLoader, nil
}
