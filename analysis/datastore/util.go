package datastore

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	"github.com/pkg/errors"
)

func RequireSingleResult(results []analysis.Artifact) (analysis.Artifact, error) {
	l := len(results)

	if l == 0 {
		return analysis.Artifact{}, errors.New("expected 1 analysis, found 0")
	} else if l > 1 {
		return analysis.Artifact{}, fmt.Errorf("expected 1 analysis, found %d", l)
	}

	return results[0], nil
}
