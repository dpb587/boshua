package datastore

import (
	"github.com/dpb587/boshua/pivnetfile"
)

func GetArtifact(index Index, f FilterParams) (pivnetfile.Artifact, error) {
	results, err := index.GetArtifacts(f)
	if err != nil {
		return pivnetfile.Artifact{}, err
	}

	l := len(results)

	if l == 0 {
		return pivnetfile.Artifact{}, NoMatchErr
	} else if l > 1 {
		return pivnetfile.Artifact{}, MultipleMatchErr
	}

	return results[0], nil
}
