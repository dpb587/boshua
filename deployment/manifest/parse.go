package manifest

import (
	"fmt"
	"strings"

	"github.com/cppforlife/go-patch/patch"
	"github.com/dpb587/boshua/osversion"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func Parse(manifestBytes []byte, localStemcell osversion.Reference) (*Manifest, error) {
	var parsed parseManifest

	err := yaml.Unmarshal(manifestBytes, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "parsing manifest")
	}

	var parsedRaw map[interface{}]interface{}

	err = yaml.Unmarshal(manifestBytes, &parsedRaw)
	if err != nil {
		return nil, errors.Wrap(err, "parsing raw manifest")
	}

	var releaseRequirements []ReleasePatch
	var stemcellRequirements []StemcellPatch

	var stemcell parseManifestStemcell

	if len(parsed.Stemcells) == 1 {
		stemcell = parsed.Stemcells[0]
	} else if len(parsed.ResourcePools) == 1 && parsed.ResourcePools[0].Stemcell != nil {
		stemcell = *parsed.ResourcePools[0].Stemcell
	} else {
		// no stemcell; nothing to do
		return &Manifest{parsed: parsedRaw}, nil
	}

	if stemcell.OS == "" {
		if strings.Contains(stemcell.Name, "-ubuntu-trusty-") {
			stemcell.OS = "ubuntu-trusty"
		} else if strings.Contains(stemcell.Name, "-ubuntu-xenial-") {
			stemcell.OS = "ubuntu-xenial"
		} else if strings.Contains(stemcell.Name, "-centos-7-") {
			stemcell.OS = "centos-7"
		} else {
			// no known os; nothing to do
			return &Manifest{parsed: parsedRaw}, nil
		}
	}

	// hacky emulation of release parsing
	// TODO support multiple stemcells
	stemcellRequirements = append(
		stemcellRequirements,
		StemcellPatch{
			OS:      stemcell.OS,
			Name:    stemcell.Name,
			Version: stemcell.Version,
		},
	)

	var cloudProviderRelease *parseManifestReleaseRef
	var cloudProviderReleaseInstalled bool

	if parsed.CloudProvider.Template != nil {
		cloudProviderRelease = parsed.CloudProvider.Template

		for _, release := range parsed.InstalledReleases() {
			if release.Release == cloudProviderRelease.Release {
				cloudProviderReleaseInstalled = true

				break
			}
		}
	}

	for releaseIdx, release := range parsed.Releases {
		if release.Version == "latest" || strings.HasSuffix(release.Version, ".latest") {
			// too dynamic for now; ignore
			continue
		} else if release.Stemcell != nil {
			// already compiled; ignore
			continue
		} else if cloudProviderRelease != nil && release.Name == cloudProviderRelease.Release {
			if !cloudProviderReleaseInstalled {
				// not installed and bosh-init can't use it; ignore for now
				continue
			} else if localStemcell.Name != stemcell.Name || localStemcell.Version != stemcell.Version {
				// used by both remote and local; ignore for now
				continue
			}
		}

		releasePatch := ReleasePatch{
			Name:    release.Name,
			Version: release.Version,
			Source: ReleasePatchRef{
				Sha1: release.Sha1,
				URL:  release.URL,
			},
			Stemcell: Stemcell{
				OS:      stemcell.OS,
				Version: stemcell.Version,
			},
			pointer: patch.MustNewPointerFromString(fmt.Sprintf("/releases/%d", releaseIdx)),
		}

		releaseRequirements = append(releaseRequirements, releasePatch)
	}

	return &Manifest{
		parsed:               parsedRaw,
		releaseRequirements:  releaseRequirements,
		stemcellRequirements: stemcellRequirements,
	}, nil
}
