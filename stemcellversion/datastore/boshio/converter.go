package boshio

import (
	"strings"

	"github.com/dpb587/boshua/stemcellversion"
)

func ConvertFileNameToReference(name string) *stemcellversion.Reference {
	ref := &stemcellversion.Reference{}

	nameSplit := strings.Split(name, "-")

	if nameSplit[0] == "light" {
		ref.Light = true
		nameSplit = nameSplit[1:]
	}

	if nameSplit[0:1] != []string{"bosh", "stemcell"} {
		// unexpected format
		return nil
	}

	nameSplit = nameSplit[2:]

	ref.Version = nameSplit[0]
	nameSplit = nameSplit[1]

	ref.IaaS = nameSplit[0]
	nameSplit = nameSplit[1:]

	ref.Hypervisor = nameSplit[0]
	nameSplit = nameSplit[1:]
	if hypervisor == "xen" && nameSplit[0] == "hvm" {
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

	if len(nameSplit) != 0 {
		// probably disk type
		return nil
	}

	return ref
}
