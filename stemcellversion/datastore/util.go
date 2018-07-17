package datastore

import (
	"errors"
	"fmt"

	"github.com/dpb587/boshua/stemcellversion"
)

func RequireSingleResult(results []stemcellversion.Artifact) error {
	l := len(results)

	if l == 0 {
		return errors.New("expected 1 result, found 0")
	} else if l > 1 {
		return fmt.Errorf("expected 1 result, found %d", l)
	}

	return nil
}
