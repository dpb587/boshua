package boshioindex

import (
	"strings"

	"github.com/dpb587/boshua/stemcellversion"
)

func ConvertFileNameToReference(name string) *stemcellversion.Artifact {
	ref := &stemcellversion.Artifact{}

	nameSplit := strings.Split(strings.TrimSuffix(name, ".tgz"), "-")

	if nameSplit[0] == "light" {
		ref.Flavor = "light"
		nameSplit = nameSplit[1:]
	} else {
		// TODO light-china?
		ref.Flavor = "heavy"
	}

	if nameSplit[0] != "bosh" || nameSplit[1] != "stemcell" {
		// unexpected format
		return nil
	}

	nameSplit = nameSplit[2:]

	ref.Version = nameSplit[0]
	nameSplit = nameSplit[1:]

	ref.IaaS = nameSplit[0]
	nameSplit = nameSplit[1:]

	ref.Hypervisor = nameSplit[0]
	nameSplit = nameSplit[1:]
	if ref.Hypervisor == "xen" && nameSplit[0] == "hvm" {
		ref.Hypervisor = strings.Join([]string{ref.Hypervisor, nameSplit[0]}, "-")
		nameSplit = nameSplit[1:]
	}

	ref.OS = nameSplit[0]
	nameSplit = nameSplit[1:]
	if !strings.HasPrefix(ref.OS, "windows") {
		ref.OS = strings.Join([]string{ref.OS, nameSplit[0]}, "-")
		nameSplit = nameSplit[1:]
	}

	if nameSplit[0] != "go_agent" {
		// undesired?
		return nil
	}

	nameSplit = nameSplit[1:]

	if len(nameSplit) > 0 {
		ref.DiskFormat = nameSplit[0]

		nameSplit = nameSplit[1:]
	}

	if len(nameSplit) != 0 {
		// dunno
		return nil
	}

	return ref
}
