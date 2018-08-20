package types

import (
	"regexp"

	"github.com/pkg/errors"
)

type Regexp struct {
	*regexp.Regexp
}

func (r *Regexp) UnmarshalYAML(data []byte) error {
	parsed, err := regexp.Compile(string(data))
	if err != nil {
		return errors.Wrap(err, "parsing regexp")
	}

	r.Regexp = parsed

	return nil
}
