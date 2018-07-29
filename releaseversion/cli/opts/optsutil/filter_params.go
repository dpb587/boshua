package opts

import (
	"fmt"

	"github.com/dpb587/boshua/releaseversion/datastore"
)

func ArgsFromFilterParams(f datastore.FilterParams) []string {
	args := []string{}

	if f.NameExpected {
		args = append(args, fmt.Sprintf("--release-name=%s", f.Name))
	}

	if f.VersionExpected {
		args = append(args, fmt.Sprintf("--release-version=%s", f.Version))
	}

	if f.ChecksumExpected {
		args = append(args, fmt.Sprintf("--release-checksum=%s", f.Checksum))
	}

	if f.URIExpected {
		args = append(args, fmt.Sprintf("--release-url=%s", f.URI))
	}

	if f.LabelsExpected {
		for _, label := range f.Labels {
			args = append(args, fmt.Sprintf("--release-label=%s", label))
		}
	}

	return args
}
