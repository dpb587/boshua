package opts

import (
	"fmt"

	"github.com/dpb587/boshua/pivnetfile/datastore"
)

func ArgsFromFilterParams(f datastore.FilterParams) []string {
	args := []string{}

	if f.ProductNameExpected {
		args = append(args, fmt.Sprintf("--pivnet-product=%s", f.ProductName))
	}

	if f.ReleaseIDExpected {
		args = append(args, fmt.Sprintf("--pivnet-release-id=%s", f.ReleaseID))
	}

	if f.FileIDExpected {
		args = append(args, fmt.Sprintf("--pivnet-file-id=%s", f.FileID))
	}

	return args
}
