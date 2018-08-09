package datastore

import (
	"github.com/dpb587/boshua/releaseversion"
)

func GetArtifact(index Index, f FilterParams) (releaseversion.Artifact, error) {
	results, err := index.GetArtifacts(f)
	if err != nil {
		return releaseversion.Artifact{}, err
	}

	l := len(results)

	if l == 0 {
		return releaseversion.Artifact{}, NoMatchErr
	} else if l > 1 {
		return releaseversion.Artifact{}, MultipleMatchErr
	}

	return results[0], nil
}
