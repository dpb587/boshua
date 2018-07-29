package datastore

import (
	"fmt"

	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/pkg/errors"
)

func RequireSingleResult(results []compilation.Artifact) (compilation.Artifact, error) {
	l := len(results)

	if l == 0 {
		return compilation.Artifact{}, errors.New("expected 1 compiled release version, found 0")
	} else if l > 1 {
		return compilation.Artifact{}, fmt.Errorf("expected 1 compiled release version, found %d", l)
	}

	return results[0], nil
}
