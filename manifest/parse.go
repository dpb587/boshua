package manifest

import (
	"fmt"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/cppforlife/go-patch/patch"
)

func Parse(manifestBytes []byte) (*Manifest, error) {
	var parsed map[interface{}]interface{}

	err := yaml.Unmarshal(manifestBytes, &parsed)
	if err != nil {
		return nil, fmt.Errorf("parsing manifest: %v", err)
	}

	var requirements []Release

	stemcell, err := parseStemcell(parsed)
	if err != nil {
		return nil, fmt.Errorf("parsing stemcell: %v", err)
	} else if stemcell.Version == "latest" || strings.HasSuffix(stemcell.Version, ".latest") {
		return &Manifest{
			parsed: parsed,
		}, nil
	}

	cloudProviderReleaseName, err := parseCloudProviderReleaseName(parsed)
	if err != nil {
		return nil, fmt.Errorf("parsing cloud provider", err)
	}

	releasesArray, ok := parsed["releases"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("releases is expected to be an array of hashes")
	}

	for releaseIdx, releaseRaw := range releasesArray {
		releaseStruct, ok := releaseRaw.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("parsing release %d: expected hash", releaseIdx)
		}

		release, err := parseRelease(releaseIdx, releaseStruct)
		if err != nil {
			return nil, fmt.Errorf("parsing release %d: %v", releaseIdx, err)
		} else if release == nil {
			// unsupported or missing fields
			continue
		} else if release.Version == "latest" || strings.HasSuffix(release.Version, ".latest") {
			continue
		} else if release.IsCompiled() {
			continue
		}

		release.Stemcell = *stemcell
		release.pointer = patch.MustNewPointerFromString(fmt.Sprintf("/releases/%d", releaseIdx))

		requirements = append(requirements, *release)
	}

	return &Manifest{
		parsed:       parsed,
		requirements: requirements,
	}, nil
}

func parseCloudProviderReleaseName(parsed map[interface{}]interface{}) (*string, error) {
	cloudProviderRaw, ok := parsed["cloud_provider"]
	if !ok {
		return nil, nil
	}

	cloudProvider, ok := cloudProviderRaw.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("cloud_provider is expected to be a hash: %#+v", cloudProviderRaw)
	}

	templateRaw, ok := cloudProvider["template"]
	if !ok {
		return nil, nil
	}

	template, ok := templateRaw.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("cloud_provider/template is expected to be a hash: %#+v", templateRaw)
	}

	releaseRaw, ok := template["release"]
	if !ok {
		return nil, nil
	}

	release, ok := releaseRaw.(string)
	if !ok {
		return nil, fmt.Errorf("cloud_provider/template/release is expected to be a string: %#+v", releaseRaw)
	}

	return &release, nil
}

func parseRelease(idx int, parsed map[interface{}]interface{}) (*Release, error) {
	release := &Release{}

	name, ok := parsed["name"]
	if !ok {
		return nil, fmt.Errorf("expected field: name")
	}

	release.Name, ok = name.(string)
	if !ok {
		return nil, fmt.Errorf("expected string field: name")
	}

	switch parsed["version"].(type) {
	case string:
		release.Version = parsed["version"].(string)
	case int:
		release.Version = strconv.Itoa(parsed["version"].(int))
	default:
		panic("expected field: version")
	}

	sha1, ok := parsed["sha1"]
	if !ok {
		return nil, nil
	}

	release.Source.Sha1, ok = sha1.(string)
	if !ok {
		return nil, fmt.Errorf("expected string field: sha1")
	}

	url, ok := parsed["url"]
	if !ok {
		return nil, fmt.Errorf("expected field: url")
	}

	release.Source.URL, ok = url.(string)
	if !ok {
		return nil, fmt.Errorf("expected string field: url")
	}

	return release, nil
}

func parseStemcell(parsed map[interface{}]interface{}) (*Stemcell, error) {
	var stemcellRaw map[interface{}]interface{}
	stemcell := &Stemcell{}

	_, hasStemcells := parsed["stemcells"]
	_, hasResourcePools := parsed["resource_pools"]

	if hasStemcells {
		stemcellsSlice, ok := parsed["stemcells"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("stemcells is expected to be an array")
		} else if len(stemcellsSlice) != 1 {
			return nil, fmt.Errorf("stemcells is expected to have a single entry but found %d", len(stemcellsSlice))
		}

		stemcellStruct, ok := stemcellsSlice[0].(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("stemcell is expected to be a hash: %#+v", stemcellsSlice[0])
		}

		stemcellRaw = stemcellStruct
	} else if hasResourcePools {
		resourcePoolsSlice, ok := parsed["resource_pools"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("resource_pools is expected to be an array")
		} else if len(resourcePoolsSlice) != 1 {
			return nil, fmt.Errorf("resource_pools is expected to have a single entry but found %d", len(resourcePoolsSlice))
		}

		resourcePoolStruct, ok := resourcePoolsSlice[0].(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("resource_pools is expected to be a hash: %#+v", resourcePoolsSlice[0])
		}

		stemcellStruct, ok := resourcePoolStruct["stemcell"].(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("resource_pools is expected to have a stemcell hash")
		}

		stemcellRaw = stemcellStruct
	} else {
		return nil, fmt.Errorf("no stemcell found")
	}

	var ok bool
	stemcellOS, okOS := stemcellRaw["os"]
	stemcellName, okName := stemcellRaw["name"]

	if okOS {
		stemcell.OS, ok = stemcellOS.(string)
		if !ok {
			return nil, fmt.Errorf("stemcell expects string field: os")
		}
	} else if okName {
		stemcellNameString, ok := stemcellName.(string)
		if !ok {
			return nil, fmt.Errorf("stemcell name is expected to be string")
		}

		if strings.Contains(stemcellNameString, "-ubuntu-trusty-") {
			stemcell.OS = "ubuntu-trusty"
		} else if strings.Contains(stemcellNameString, "-ubuntu-xenial-") {
			stemcell.OS = "ubuntu-xenial"
		} else if strings.Contains(stemcellNameString, "-centos-7-") {
			stemcell.OS = "centos-7"
		}
	} else {
		return nil, fmt.Errorf("stemcell expects field: os or name")
	}

	stemcellVersion, ok := stemcellRaw["version"]
	if !ok {
		return nil, fmt.Errorf("stemcell expects field: version")
	}

	stemcell.Version, ok = stemcellVersion.(string)
	if !ok {
		return nil, fmt.Errorf("stemcell expects string field: version")
	}

	return stemcell, nil
}
