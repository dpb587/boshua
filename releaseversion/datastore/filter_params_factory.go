package datastore

import (
	"errors"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/util/semverutil"
)

func FilterParamsFromSlug(slug string) FilterParams {
	split := strings.SplitN(slug, "/", 2)

	f := FilterParams{
		NameExpected: true,
		Name:         split[0],
	}

	if len(split) == 2 {
		f.VersionExpected = true
		f.Version = split[1]
	}

	return f
}

func FilterParamsFromArtifact(artifact releaseversion.Artifact) FilterParams {
	f := FilterParams{
		NameExpected: true,
		Name:         artifact.Name,

		VersionExpected: true,
		Version:         artifact.Version,

		// Labels are relative/subjective; irrelevant to artifact identity
	}

	// TODO should be tracking original checksum/uri?
	if len(artifact.MetalinkFile().Hashes) > 0 {
		f.ChecksumExpected = true
		f.Checksum = metalinkutil.HashToChecksum(metalinkutil.PreferredHash(artifact.MetalinkFile().Hashes)).String()
	}

	return f
}

func FilterParamsFromReference(ref releaseversion.Reference) FilterParams {
	f := FilterParams{
		NameExpected: true,
		Name:         ref.Name,

		VersionExpected: true,
		Version:         ref.Version,

		// Labels are relative/subjective; irrelevant to artifact identity
	}

	if len(ref.Checksums) > 0 {
		// TODO use strongest
		f.ChecksumExpected = true
		f.Checksum = ref.Checksums[0].String()
	}

	// TODO uri

	return f
}

func FilterParamsFromMap(args map[string]interface{}) (FilterParams, error) {
	f := FilterParams{}

	f.Name, f.NameExpected = args["name"].(string)
	f.Version, f.VersionExpected = args["version"].(string)
	f.Checksum, f.ChecksumExpected = args["checksum"].(string)
	f.URI, f.URIExpected = args["uri"].(string)

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
