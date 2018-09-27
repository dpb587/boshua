package provider

import (
	"github.com/pkg/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/boshua/metalink/file/metaurl/boshreleasesource"
	"github.com/dpb587/metalink/file/url"
	fileurl "github.com/dpb587/metalink/file/url/file"
	ftpurl "github.com/dpb587/metalink/file/url/ftp"
	httpurl "github.com/dpb587/metalink/file/url/http"
	s3url "github.com/dpb587/metalink/file/url/s3"
	"github.com/dpb587/metalink/file/metaurl"
	downloaderurl "github.com/dpb587/boshua/artifact/downloader/url"
	urlfilteredloader "github.com/dpb587/boshua/metalink/file/url/filteredloader"
)

func (c *Config) SetDownloaderURLFactory(f downloaderurl.Factory) {
	c.downloaderUrlFactory = f
}

func (c *Config) GetDownloader() (url.Loader, metaurl.Loader, error) {
	urlLoader := url.NewLoaderFactory()

	for _, cfg := range c.Downloaders.URLHandlers {
		handler, err := c.downloaderUrlFactory.Create(downloaderurl.ProviderName(cfg.Type), cfg.Name, cfg.Options)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "creating url handler (name: %s)", cfg.Name)
		}

		if len(cfg.Include) > 0 || len(cfg.Exclude) > 0 {
			handler = urlfilteredloader.NewLoader(handler, cfg.Include.AsRegexp(), cfg.Exclude.AsRegexp())
		}

		urlLoader.Add(handler)
	}

	metaurlLoader := metaurl.NewLoaderFactory()

	if !c.Downloaders.DisableDefaultHandlers {
		logger := boshlog.NewLogger(boshlog.LevelError)
		fs := boshsys.NewOsFileSystem(logger)
		file := fileurl.NewLoader(fs)

		urlLoader.Add(ftpurl.Loader{})
		urlLoader.Add(httpurl.Loader{})
		urlLoader.Add(s3url.NewLoader("", ""))

		// TODO avoid file access by default?
		urlLoader.Add(file)
		urlLoader.Add(fileurl.NewEmptyLoader(file))

		metaurlLoader.Add(boshreleasesource.Loader{})
	}

	return urlLoader, metaurlLoader, nil
}
