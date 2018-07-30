package graphql

import (
	"strings"

	"github.com/dpb587/boshua/releaseversion/datastore"
)

func BuildListQueryArgs(f datastore.FilterParams) (string, string, map[string]interface{}) {
	var queryFilter, queryVarsTypes []string
	var queryVars = map[string]interface{}{}

	if f.NameExpected {
		queryFilter = append(queryFilter, "name: $qReleaseName")
		queryVarsTypes = append(queryVarsTypes, "$qReleaseName: String!")
		queryVars["qReleaseName"] = f.Name
	}

	if f.VersionExpected {
		queryFilter = append(queryFilter, "version: $qReleaseVersion")
		queryVarsTypes = append(queryVarsTypes, "$qReleaseVersion: String!")
		queryVars["qReleaseVersion"] = f.Version
	}

	if f.ChecksumExpected {
		queryFilter = append(queryFilter, "checksum: $qReleaseChecksum")
		queryVarsTypes = append(queryVarsTypes, "$qReleaseChecksum: String!")
		queryVars["qReleaseChecksum"] = f.Checksum
	}

	if f.URIExpected {
		queryFilter = append(queryFilter, "uri: $qReleaseUri")
		queryVarsTypes = append(queryVarsTypes, "$qReleaseUri: String!")
		queryVars["qReleaseUri"] = f.URI
	}

	if f.LabelsExpected {
		queryFilter = append(queryFilter, "labels: $qReleaseLabels")
		queryVarsTypes = append(queryVarsTypes, "$qReleaseLabels: [String!]")
		queryVars["qReleaseLabels"] = f.Labels
	}

	return strings.Join(queryFilter, ", "), strings.Join(queryVarsTypes, ", "), queryVars
}
