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

type RegexpList []*Regexp

func (rl RegexpList) AsRegexp() []*regexp.Regexp {
	var as []*regexp.Regexp

	for _, r := range rl {
		as = append(as, r.Regexp)
	}

	return as
}
