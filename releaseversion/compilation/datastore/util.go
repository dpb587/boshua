package datastore

import (
	"github.com/dpb587/boshua/releaseversion/compilation"
)

func GetCompilationArtifact(index Index, f FilterParams) (compilation.Artifact, error) {
	results, err := index.GetCompilationArtifacts(f)
	if err != nil {
		return compilation.Artifact{}, err
	}

	l := len(results)

	if l == 0 {
		return compilation.Artifact{}, NoMatchErr
	} else if l > 1 {
		return compilation.Artifact{}, MultipleMatchErr
	}

	return results[0], nil
}
