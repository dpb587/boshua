package graphql

import (
	"strings"

	"github.com/dpb587/boshua/releaseversion/datastore"
)

func BuildListQueryArgs(f datastore.FilterParams) (string, string, map[string]interface{}) {
	var queryFilter, queryVarsTypes []string
	var queryVars = map[string]interface{}{}

	if f.NameExpected {
		queryFilter = append(queryFilter, "name: $queryName")
		queryVarsTypes = append(queryVarsTypes, "$queryName: String!")
		queryVars["queryName"] = f.Name
	}

	if f.VersionExpected {
		queryFilter = append(queryFilter, "version: $queryVersion")
		queryVarsTypes = append(queryVarsTypes, "$queryVersion: String!")
		queryVars["queryVersion"] = f.Version
	}

	if f.ChecksumExpected {
		queryFilter = append(queryFilter, "checksum: $queryChecksum")
		queryVarsTypes = append(queryVarsTypes, "$queryChecksum: String!")
		queryVars["queryChecksum"] = f.Checksum
	}

	if f.URIExpected {
		queryFilter = append(queryFilter, "uri: $queryURI")
		queryVarsTypes = append(queryVarsTypes, "$queryURI: String!")
		queryVars["queryURI"] = f.URI
	}

	if f.LabelsExpected {
		queryFilter = append(queryFilter, "labels: $queryLabels")
		queryVarsTypes = append(queryVarsTypes, "$queryLabels: [String!]")
		queryVars["queryLabels"] = f.Labels
	}

	return strings.Join(queryFilter, ", "), strings.Join(queryVarsTypes, ", "), queryVars
}
