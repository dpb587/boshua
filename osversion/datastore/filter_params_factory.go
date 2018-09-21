package datastore

import (
	"github.com/Masterminds/semver"
	"github.com/dpb587/boshua/util/semverutil"
)

func FilterParamsFromMap(args map[string]interface{}) (FilterParams, error) {
	f := FilterParams{}

	// TODO consolidate os vs name
	if _, found := args["os"]; found {
		f.Name, f.NameExpected = args["os"].(string)
	} else {
		f.Name, f.NameExpected = args["name"].(string)
	}

	f.Version, f.VersionExpected = args["version"].(string)

	if f.VersionExpected && semverutil.IsConstraint(f.Version) {
		// ignoring errors since it can fallback to literal match
		f.VersionConstraint, _ = semver.NewConstraint(f.Version)
	}

	return f, nil
}
