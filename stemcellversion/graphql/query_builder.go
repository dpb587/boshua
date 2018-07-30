package graphql

import (
	"strings"

	"github.com/dpb587/boshua/stemcellversion/datastore"
)

func BuildListQueryArgs(f datastore.FilterParams) (string, string, map[string]interface{}) {
	var queryFilter, queryVarsTypes []string
	var queryVars = map[string]interface{}{}

	if f.OSExpected {
		queryFilter = append(queryFilter, "os: $qStemcellOs")
		queryVarsTypes = append(queryVarsTypes, "$qStemcellOs: String!")
		queryVars["qStemcellOs"] = f.OS
	}

	if f.VersionExpected {
		queryFilter = append(queryFilter, "version: $qStemcellVersion")
		queryVarsTypes = append(queryVarsTypes, "$qStemcellVersion: String!")
		queryVars["qStemcellVersion"] = f.Version
	}

	if f.IaaSExpected {
		queryFilter = append(queryFilter, "iaas: $qStemcellIaaS")
		queryVarsTypes = append(queryVarsTypes, "$qStemcellIaaS: String!")
		queryVars["qStemcellIaaS"] = f.IaaS
	}

	if f.HypervisorExpected {
		queryFilter = append(queryFilter, "hypervisor: $qStemcellHypervisor")
		queryVarsTypes = append(queryVarsTypes, "$qStemcellHypervisor: String!")
		queryVars["qStemcellHypervisor"] = f.Hypervisor
	}

	if f.DiskFormatExpected {
		queryFilter = append(queryFilter, "diskFormat: $qStemcellDiskFormat")
		queryVarsTypes = append(queryVarsTypes, "$qStemcellDiskFormat: String!")
		queryVars["qStemcellDiskFormat"] = f.DiskFormat
	}

	if f.FlavorExpected {
		queryFilter = append(queryFilter, "flavor: $qStemcellFlavor")
		queryVarsTypes = append(queryVarsTypes, "$qStemcellFlavor: String!")
		queryVars["qStemcellFlavor"] = f.Flavor
	}

	if f.LabelsExpected {
		queryFilter = append(queryFilter, "labels: $qStemcellLabels")
		queryVarsTypes = append(queryVarsTypes, "$qStemcellLabels: [String!]")
		queryVars["qStemcellLabels"] = f.Labels
	}

	return strings.Join(queryFilter, ", "), strings.Join(queryVarsTypes, ", "), queryVars
}
