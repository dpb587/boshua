package datastore

import (
	"github.com/dpb587/boshua/analysis"
)

func GetAnalysisArtifact(index Index, ref analysis.Reference) (analysis.Artifact, error) {
	results, err := index.GetAnalysisArtifacts(ref)
	if err != nil {
		return analysis.Artifact{}, err
	}

	l := len(results)

	if l == 0 {
		return analysis.Artifact{}, NoMatchErr
	} else if l > 1 {
		return analysis.Artifact{}, MultipleMatchErr
	}

	return results[0], nil
}
