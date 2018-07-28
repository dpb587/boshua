package datastore

import (
	"github.com/Masterminds/semver"
)

type FilterParams struct {
	OSExpected bool
	OS         string

	VersionExpected   bool
	Version           string
	VersionConstraint *semver.Constraints

	IaaSExpected bool
	IaaS         string

	HypervisorExpected bool
	Hypervisor         string

	DiskFormatExpected bool
	DiskFormat         string

	FlavorExpected bool
	Flavor         string

	LabelsExpected bool // TODO unnecessary? implied by len > 0
	Labels         []string
}

func (f *FilterParams) OSSatisfied(actual string) bool {
	if !f.OSExpected {
		return true
	}

	return f.OS == actual
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

func (f *FilterParams) IaaSSatisfied(actual string) bool {
	if !f.IaaSExpected {
		return true
	}

	return f.IaaS == actual
}

func (f *FilterParams) HypervisorSatisfied(actual string) bool {
	if !f.HypervisorExpected {
		return true
	}

	return f.Hypervisor == actual
}

func (f *FilterParams) DiskFormatSatisfied(actual string) bool {
	if !f.DiskFormatExpected {
		return true
	}

	return f.DiskFormat == actual
}

func (f *FilterParams) FlavorSatisfied(actual string) bool {
	if !f.FlavorExpected {
		return true
	}

	return f.Flavor == actual
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
