package graphql

import (
	"strings"

	"github.com/dpb587/boshua/stemcellversion/datastore"
)

func BuildListQueryArgs(f datastore.FilterParams, l datastore.LimitParams) (string, string, map[string]interface{}) {
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

	if l.MinExpected {
		queryFilter = append(queryFilter, "limitMin: $qStemcellLimitMin")
		queryVarsTypes = append(queryVarsTypes, "$qStemcellLimitMin: Float!")
		queryVars["qStemcellLimitMin"] = l.Min
	}

	if l.MaxExpected {
		queryFilter = append(queryFilter, "limitMax: $qStemcellLimitMax")
		queryVarsTypes = append(queryVarsTypes, "$qStemcellLimitMax: Float!")
		queryVars["qStemcellLimitMax"] = l.Max
	}

	if l.LimitExpected {
		queryFilter = append(queryFilter, "limitFirst: $qStemcellLimitFirst")
		queryVarsTypes = append(queryVarsTypes, "$qStemcellLimitFirst: Float!")
		queryVars["qStemcellLimitFirst"] = l.Limit
	}

	if l.OffsetExpected {
		queryFilter = append(queryFilter, "limitOffset: $qStemcellLimitOffset")
		queryVarsTypes = append(queryVarsTypes, "$qStemcellLimitOffset: Float!")
		queryVars["qStemcellLimitOffset"] = l.Offset
	}

	return strings.Join(queryFilter, ", "), strings.Join(queryVarsTypes, ", "), queryVars
}
