package datastore

import (
	"github.com/dpb587/boshua/stemcellversion"
)

func GetArtifact(index Index, f FilterParams) (stemcellversion.Artifact, error) {
	results, err := index.GetArtifacts(f)
	if err != nil {
		return stemcellversion.Artifact{}, err
	}

	l := len(results)

	if l == 0 {
		return stemcellversion.Artifact{}, NoMatchErr
	} else if l > 1 {
		return stemcellversion.Artifact{}, MultipleMatchErr
	}

	return results[0], nil
}
