package manifest

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/dpb587/boshua/osversion"

	"github.com/cppforlife/go-patch/patch"
)

func Parse(manifestBytes []byte, localStemcell osversion.Reference) (*Manifest, error) {
	var parsed parseManifest

	err := yaml.Unmarshal(manifestBytes, &parsed)
	if err != nil {
		return nil, fmt.Errorf("parsing manifest: %v", err)
	}

	var parsedRaw map[interface{}]interface{}

	err = yaml.Unmarshal(manifestBytes, &parsedRaw)
	if err != nil {
		return nil, fmt.Errorf("parsing raw manifest: %v", err)
	}

	var requirements []ReleasePatch

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
		} else if cloudProviderReleaseInstalled && release.Name == cloudProviderRelease.Release {
			if localStemcell.Name != stemcell.Name || localStemcell.Version != stemcell.Version {
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

		requirements = append(requirements, releasePatch)
	}

	return &Manifest{
		parsed:       parsedRaw,
		requirements: requirements,
	}, nil
}
