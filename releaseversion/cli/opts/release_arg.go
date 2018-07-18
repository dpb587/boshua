package opts

import (
	"fmt"
	"strings"
)

type Release struct {
	Name    string
	Version string
}

func (r Release) String() string {
	return fmt.Sprintf("%s/%s", r.Name, r.Version)
}

func (r *Release) UnmarshalFlag(data string) error {
	split := strings.Split(data, "/")

	if len(split) != 2 {
		return fmt.Errorf("expected release format of name/version: %s", data)
	}

	r.Name = split[0]
	r.Version = split[1]

	return nil
}
