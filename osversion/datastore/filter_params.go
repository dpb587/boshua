package datastore

import (
	"github.com/Masterminds/semver"
)

type FilterParams struct {
	NameExpected bool
	Name         string

	VersionExpected   bool
	Version           string
	VersionConstraint *semver.Constraints
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
