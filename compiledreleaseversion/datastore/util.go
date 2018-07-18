package datastore

import (
	"errors"
	"fmt"

	"github.com/dpb587/boshua/compiledreleaseversion"
)

func RequireSingleResult(results []compiledreleaseversion.Artifact) (compiledreleaseversion.Artifact, error) {
	l := len(results)

	if l == 0 {
		return compiledreleaseversion.Artifact{}, errors.New("expected 1 compiled release version, found 0")
	} else if l > 1 {
		return compiledreleaseversion.Artifact{}, fmt.Errorf("expected 1 compiled release version, found %d", l)
	}

	return results[0], nil
}
