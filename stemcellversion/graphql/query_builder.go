package graphql

import (
	"strings"

	"github.com/dpb587/boshua/stemcellversion/datastore"
)

func BuildListQueryArgs(f datastore.FilterParams) (string, string, map[string]interface{}) {
	var queryFilter, queryVarsTypes []string
	var queryVars = map[string]interface{}{}

	if f.OSExpected {
		queryFilter = append(queryFilter, "os: $queryOS")
		queryVarsTypes = append(queryVarsTypes, "$queryOS: String!")
		queryVars["queryOS"] = f.OS
	}

	if f.VersionExpected {
		queryFilter = append(queryFilter, "version: $queryVersion")
		queryVarsTypes = append(queryVarsTypes, "$queryVersion: String!")
		queryVars["queryVersion"] = f.Version
	}

	if f.IaaSExpected {
		queryFilter = append(queryFilter, "iaas: $queryIaaS")
		queryVarsTypes = append(queryVarsTypes, "$queryIaaS: String!")
		queryVars["queryIaaS"] = f.IaaS
	}

	if f.HypervisorExpected {
		queryFilter = append(queryFilter, "hypervisor: $queryHypervisor")
		queryVarsTypes = append(queryVarsTypes, "$queryHypervisor: String!")
		queryVars["queryHypervisor"] = f.Hypervisor
	}

	if f.DiskFormatExpected {
		queryFilter = append(queryFilter, "diskFormat: $queryDiskFormat")
		queryVarsTypes = append(queryVarsTypes, "$queryDiskFormat: String!")
		queryVars["queryDiskFormat"] = f.DiskFormat
	}

	if f.FlavorExpected {
		queryFilter = append(queryFilter, "flavor: $queryFlavor")
		queryVarsTypes = append(queryVarsTypes, "$queryFlavor: String!")
		queryVars["queryFlavor"] = f.Flavor
	}

	if f.LabelsExpected {
		queryFilter = append(queryFilter, "labels: $queryLabels")
		queryVarsTypes = append(queryVarsTypes, "$queryLabels: [String!]")
		queryVars["queryLabels"] = f.Labels
	}

	return strings.Join(queryFilter, ", "), strings.Join(queryVarsTypes, ", "), queryVars
}
