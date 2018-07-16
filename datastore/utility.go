package datastore

import (
	"errors"
	"fmt"
)

func RequireSingleResult(results interface{}) error {
	l := len(results.([]interface{}))

	if l == 0 {
		return errors.New("expected 1 result, found 0")
	} else if l > 1 {
		return fmt.Errorf("expected 1 result, found %d", l)
	}

	return nil
}
