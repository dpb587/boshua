package datastore

import (
	"errors"
	"fmt"

	"github.com/dpb587/boshua/releaseversion"
)

func RequireSingleResult(results []releaseversion.Artifact) (releaseversion.Artifact, error) {
	l := len(results)

	if l == 0 {
		return releaseversion.Artifact{}, errors.New("expected 1 release version, found 0")
	} else if l > 1 {
		return releaseversion.Artifact{}, fmt.Errorf("expected 1 release version, found %d", l)
	}

	return results[0], nil
}
