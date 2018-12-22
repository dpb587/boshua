package datastore

import (
	"errors"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/util/semverutil"
)

// TODO should not be panicing
func FilterParamsFromSlug(slug string) FilterParams {
	f := FilterParams{}

	slugSplit := strings.SplitN(slug, "/", 2)
	nameSplit := strings.Split(slugSplit[0], "-")

	if nameSplit[0] == "light" {
		f.FlavorExpected = true
		f.Flavor = "light"

		nameSplit = nameSplit[1:]
	} else {
		// TODO light-china?
		f.FlavorExpected = true
		f.Flavor = "heavy"
	}

	if nameSplit[0] != "bosh" {
		// unexpected format
		panic("TODO") // TODO !panic
	}

	nameSplit = nameSplit[1:]

	f.IaaSExpected = true
	f.IaaS = nameSplit[0]
	nameSplit = nameSplit[1:]

	f.HypervisorExpected = true
	f.Hypervisor = nameSplit[0]
	nameSplit = nameSplit[1:]

	if f.Hypervisor == "xen" && nameSplit[0] == "hvm" {
		f.Hypervisor = strings.Join([]string{f.Hypervisor, nameSplit[0]}, "-")
		nameSplit = nameSplit[1:]
	}

	f.OSExpected = true
	f.OS = nameSplit[0]
	nameSplit = nameSplit[1:]

	if !strings.HasPrefix(f.OS, "windows") {
		f.OS = strings.Join([]string{f.OS, nameSplit[0]}, "-")
		nameSplit = nameSplit[1:]
	}

	if nameSplit[0] != "go_agent" {
		// undesired?
		panic("TODO") // TODO !panic
	}

	nameSplit = nameSplit[1:]

	if len(nameSplit) != 0 {
		// probably disk type
		panic("TODO") // TODO !panic
	}

	if len(slugSplit) > 1 {
		f.VersionExpected = true
		f.Version = slugSplit[1]
	}

	return f
}

func FilterParamsFromArtifact(artifact stemcellversion.Artifact) FilterParams {
	return FilterParams{
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

		// Labels are relative/subjective; irrelevant to artifact identity
	}
}

func FilterParamsFromReference(ref stemcellversion.Reference) FilterParams {
	return FilterParams{
		OSExpected: true,
		OS:         ref.OS,

		VersionExpected: true,
		Version:         ref.Version,

		IaaSExpected: true,
		IaaS:         ref.IaaS,

		HypervisorExpected: true,
		Hypervisor:         ref.Hypervisor,

		DiskFormatExpected: true,
		DiskFormat:         ref.DiskFormat,

		FlavorExpected: true,
		Flavor:         ref.Flavor,

		// Labels are relative/subjective; irrelevant to artifact identity
	}
}

func FilterParamsFromMap(args map[string]interface{}) (FilterParams, error) {
	f := FilterParams{}

	f.OS, f.OSExpected = args["os"].(string)
	f.Version, f.VersionExpected = args["version"].(string)
	f.IaaS, f.IaaSExpected = args["iaas"].(string)
	f.Hypervisor, f.HypervisorExpected = args["hypervisor"].(string)
	f.Flavor, f.FlavorExpected = args["flavor"].(string)

	if f.VersionExpected && semverutil.IsConstraint(f.Version) {
		// ignoring errors since it can fallback to literal match
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
