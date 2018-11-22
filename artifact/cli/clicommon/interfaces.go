package clicommon

import (
	"github.com/dpb587/metalink/transfer"
)

type DownloaderGetter func () (transfer.Transfer, error)
