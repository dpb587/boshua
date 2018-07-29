package optsutil

import (
	"fmt"

	"github.com/dpb587/boshua/stemcellversion/datastore"
)

func ArgsFromFilterParams(f datastore.FilterParams) []string {
	var args []string

	if f.OSExpected {
		args = append(args, fmt.Sprintf("--stemcell-os=%s", f.OS))
	}

	if f.VersionExpected {
		args = append(args, fmt.Sprintf("--stemcell-version=%s", f.Version))
	}

	if f.IaaSExpected {
		args = append(args, fmt.Sprintf("--stemcell-iaas=%s", f.IaaS))
	}

	if f.HypervisorExpected {
		args = append(args, fmt.Sprintf("--stemcell-hypervisor=%s", f.Hypervisor))
	}

	if f.FlavorExpected {
		args = append(args, fmt.Sprintf("--stemcell-flavor=%s", f.Flavor))
	}

	if f.LabelsExpected {
		for _, label := range f.Labels {
			args = append(args, fmt.Sprintf("--stemcell-label=%s", label))
		}
	}

	return args
}
