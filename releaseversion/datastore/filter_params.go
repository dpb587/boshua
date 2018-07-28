package datastore

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/dpb587/metalink"
)

type FilterParams struct {
	NameExpected bool
	Name         string

	VersionExpected   bool
	Version           string
	VersionConstraint *semver.Constraints

	ChecksumExpected bool
	Checksum         string

	URIExpected bool
	URI         string

	LabelsExpected bool // TODO unnecessary? implied by len > 0
	Labels         []string
}

func FilterParamsFromMap(args map[string]interface{}) (FilterParams, error) {
	f := FilterParams{}

	f.Name, f.NameExpected = args["name"].(string)
	f.Version, f.VersionExpected = args["version"].(string)
	f.Checksum, f.ChecksumExpected = args["checksum"].(string)
	f.URI, f.URIExpected = args["uri"].(string)

	if f.VersionExpected {
		f.VersionConstraint, _ = semver.NewConstraint(f.Version)
	}

	var labels []interface{}
	labels, f.LabelsExpected = args["labels"].([]interface{})
	for _, label := range labels {
		labelStr, ok := label.(string)
		if !ok {
			return FilterParams{}, errors.New("label: expected string")
		}

		f.Labels = append(f.Labels, labelStr)
	}

	return f, nil
}

func (f *FilterParams) NameSatisfied(actual string) bool {
	if !f.NameExpected {
		return true
	}

	return f.Name == actual
}

func (f *FilterParams) VersionSatisfied(actual string) bool {
	if !f.VersionExpected {
		return true
	} else if f.Version == actual {
		return true
	} else if f.VersionConstraint == nil {
		return false
	}

	actualVersion, err := semver.NewVersion(actual)
	if err != nil {
		return false
	}

	return f.VersionConstraint.Check(actualVersion)
}

func (f *FilterParams) ChecksumSatisfied(actual []metalink.Hash) bool {
	if !f.ChecksumExpected {
		return true
	}

	for _, hash := range actual {
		if f.Checksum == fmt.Sprintf("%s:%s", strings.Replace(hash.Type, "-", "", 1), hash.Hash) {
			return true
		}
	}

	return false
}

func (f *FilterParams) URISatisfied(actualURL []metalink.URL, actualMetaURL []metalink.MetaURL) bool {
	if !f.URIExpected {
		return true
	}

	for _, url := range actualURL {
		if f.URI == url.URL {
			return true
		}
	}

	for _, metaurl := range actualMetaURL {
		if f.URI == metaurl.URL {
			return true
		}
	}

	return false
}

func (f *FilterParams) LabelsSatisfied(actuals []string) bool {
	if !f.LabelsExpected {
		return true
	}

	for _, label := range f.Labels {
		var found bool

		for _, actual := range actuals {
			if actual == label {
				found = true

				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}
