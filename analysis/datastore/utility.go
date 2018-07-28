package datastore

import (
	"github.com/dpb587/boshua/analysis"
)

func FilterForOne(index Index, ref analysis.Reference) (analysis.Artifact, error) {
	results, err := index.GetAnalysisArtifacts(ref)
	if err != nil {
		return analysis.Artifact{}, err
	} else if len(results) == 0 {
		return analysis.Artifact{}, NoMatchErr
	} else if len(results) > 1 {
		return analysis.Artifact{}, MultipleMatchErr
	}

	return results[0], nil
}
