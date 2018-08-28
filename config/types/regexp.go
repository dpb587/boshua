package types

import (
	"regexp"

	"github.com/pkg/errors"
)

type Regexp struct {
	*regexp.Regexp
}

func (r *Regexp) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string

	err := unmarshal(&s)
	if err != nil {
		return err
	}

	parsed, err := regexp.Compile(s)
	if err != nil {
		return errors.Wrap(err, "parsing regexp")
	}

	r.Regexp = parsed

	return nil
}
