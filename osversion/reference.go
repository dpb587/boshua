package osversion

import (
	"github.com/Masterminds/semver"
)

type Reference struct {
	Name    string `json:"name"`
	Version string `json:"version"`

	semver       *semver.Version
	semverParsed bool
}

func (s Reference) Semver() *semver.Version {
	if s.semverParsed {
		return s.semver
	}

	semver, err := semver.NewVersion(s.Version)
	if err == nil {
		s.semver = semver
	}

	s.semverParsed = true

	return s.semver
}
