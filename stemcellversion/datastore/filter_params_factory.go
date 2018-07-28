package datastore

import (
	"errors"

	"github.com/Masterminds/semver"
	"github.com/dpb587/boshua/stemcellversion"
)

func FilterParamsFromArtifact(artifact stemcellversion.Artifact) *FilterParams {
	return &FilterParams{
		OSExpected: true,
		OS:         artifact.OS,

		VersionExpected: true,
		Version:         artifact.Version,

		IaaSExpected: true,
		IaaS:         artifact.IaaS,

		HypervisorExpected: true,
		Hypervisor:         artifact.Hypervisor,

		DiskFormatExpected: true,
		DiskFormat:         artifact.DiskFormat,

		FlavorExpected: true,
		Flavor:         artifact.Flavor,
	}
}

func FilterParamsFromMap(args map[string]interface{}) (*FilterParams, error) {
	f := &FilterParams{}

	f.OS, f.OSExpected = args["os"].(string)
	f.Version, f.VersionExpected = args["version"].(string)
	f.IaaS, f.IaaSExpected = args["iaas"].(string)
	f.Hypervisor, f.HypervisorExpected = args["hypervisor"].(string)
	f.Flavor, f.FlavorExpected = args["flavor"].(string)

	if f.VersionExpected {
		f.VersionConstraint, _ = semver.NewConstraint(f.Version)
	}

	var labels []interface{}
	labels, f.LabelsExpected = args["labels"].([]interface{})
	for _, label := range labels {
		labelStr, ok := label.(string)
		if !ok {
			return nil, errors.New("label: expected string")
		}

		f.Labels = append(f.Labels, labelStr)
	}

	return f, nil
}
