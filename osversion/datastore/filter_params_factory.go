package datastore

import (
	"net/url"

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

// TODO consolidate os vs stemcell
func FilterParamsFromURLValues(uv url.Values) (FilterParams, error) {
	f := FilterParams{}

	// TODO consolidate os vs name
	if values, found := uv["stemcell-os"]; found {
		// TODO validate len == 1
		f.NameExpected = true
		f.Name = values[0]
	}

	if values, found := uv["stemcell-version"]; found {
		// TODO validate len == 1
		f.VersionExpected = true
		f.Version = values[0]
	}

	if f.VersionExpected && semverutil.IsConstraint(f.Version) {
		// ignoring errors since it can fallback to literal match
		f.VersionConstraint, _ = semver.NewConstraint(f.Version)
	}

	return f, nil
}
