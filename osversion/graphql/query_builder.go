package graphql

import (
	"strings"

	"github.com/dpb587/boshua/osversion/datastore"
)

func BuildListQueryArgs(f datastore.FilterParams) (string, string, map[string]interface{}) {
	var queryFilter, queryVarsTypes []string
	var queryVars = map[string]interface{}{}

	if f.NameExpected {
		queryFilter = append(queryFilter, "os: $qOsName")
		queryVarsTypes = append(queryVarsTypes, "$qOsName: String!")
		queryVars["qOsName"] = f.Name
	}

	if f.VersionExpected {
		queryFilter = append(queryFilter, "version: $qOsVersion")
		queryVarsTypes = append(queryVarsTypes, "$qOsVersion: String!")
		queryVars["qOsVersion"] = f.Version
	}

	return strings.Join(queryFilter, ", "), strings.Join(queryVarsTypes, ", "), queryVars
}
