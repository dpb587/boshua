package args

import (
	"fmt"
	"strings"
)

type OS struct {
	Name    string
	Version string
}

func (s OS) String() string {
	return fmt.Sprintf("%s/%s", s.Name, s.Version)
}

func (s *OS) UnmarshalFlag(data string) error {
	split := strings.Split(data, "/")

	if len(split) != 2 {
		return fmt.Errorf("expected stemcell format of os/version: %s", data)
	}

	s.Name = split[0]
	s.Version = split[1]

	return nil
}
