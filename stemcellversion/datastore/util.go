package datastore

import (
	"errors"
	"fmt"

	"github.com/dpb587/boshua/stemcellversion"
)

func RequireSingleResult(results []stemcellversion.Artifact) (stemcellversion.Artifact, error) {
	l := len(results)

	if l == 0 {
		return stemcellversion.Artifact{}, errors.New("expected 1 stemcell version, found 0")
	} else if l > 1 {
		return stemcellversion.Artifact{}, fmt.Errorf("expected 1 stemcell version, found %d", l)
	}

	return results[0], nil
}
