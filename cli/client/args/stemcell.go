package args

import (
	"fmt"
	"strings"
)

type Stemcell struct {
	OS      string
	Version string
}

func (s Stemcell) String() string {
	return fmt.Sprintf("%s/%s", s.OS, s.Version)
}

func (s *Stemcell) UnmarshalFlag(data string) error {
	split := strings.Split(data, "/")

	if len(split) != 2 {
		return fmt.Errorf("expected stemcell format of os/version: %s", data)
	}

	s.OS = split[0]
	s.Version = split[1]

	return nil
}
