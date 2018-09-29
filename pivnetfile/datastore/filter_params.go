package datastore

import (
	"github.com/Masterminds/semver"
)

type FilterParams struct {
	ProductNameExpected bool
	ProductName         string

	ReleaseIDExpected bool
	ReleaseID         int

	ReleaseVersionExpected   bool
	ReleaseVersion           string
	ReleaseVersionConstraint *semver.Constraints

	FileIDExpected bool
	FileID         int

	FileNameExpected bool
	FileName         string
}

func (f *FilterParams) ProductNameSatisfied(actual string) bool {
	if !f.ProductNameExpected {
		return true
	}

	return f.ProductName == actual
}

func (f *FilterParams) ReleaseIDSatisfied(actual int) bool {
	if !f.ReleaseIDExpected {
		return true
	}

	return f.ReleaseID == actual
}

func (f *FilterParams) ReleaseVersionSatisfied(actual string) bool {
	if !f.ReleaseVersionExpected {
		return true
	} else if f.ReleaseVersion == actual {
		return true
	} else if f.ReleaseVersionConstraint == nil {
		return false
	}

	actualVersion, err := semver.NewVersion(actual)
	if err != nil {
		return false
	}

	return f.ReleaseVersionConstraint.Check(actualVersion)
}

func (f *FilterParams) FileIDSatisfied(actual int) bool {
	if !f.FileIDExpected {
		return true
	}

	return f.FileID == actual
}

func (f *FilterParams) FileNameSatisfied(actual string) bool {
	if !f.FileNameExpected {
		return true
	}

	return f.FileName == actual
}
