package clicommon

import (
	"github.com/dpb587/metalink/file/url"
	"github.com/dpb587/metalink/file/metaurl"
)

type DownloaderGetter func () (url.Loader, metaurl.Loader, error)
