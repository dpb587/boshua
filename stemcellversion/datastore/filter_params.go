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
}

func FilterParamsFromMap(args map[string]interface{}) (*FilterParams, error) {
	f := &FilterParams{}

	f.OS, f.OSExpected = args["os"].(string)
	f.Version, f.VersionExpected = args["version"].(string)
	f.IaaS, f.IaaSExpected = args["iaas"].(string)
	f.Hypervisor, f.HypervisorExpected = args["hypervisor"].(string)
	f.Flavor, f.FlavorExpected = args["light"].(string)

	if f.VersionExpected {
		f.VersionConstraint, _ = semver.NewConstraint(f.Version)
	}

	return f, nil
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
