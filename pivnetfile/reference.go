package pivnetfile

import (
	"github.com/Masterminds/semver"
)

type Reference struct {
	ProductName    string `json:"product_name"`
	ReleaseID      int    `json:"release_id"`
	FileID         int    `json:"file_id"`

	ReleaseVersion string `json:"version"`
	FileName       string `json:"file_name"`

	semver       *semver.Version
	semverParsed bool
}

func (s Reference) Semver() *semver.Version {
	if s.semverParsed {
		return s.semver
	}

	semver, err := semver.NewVersion(s.ReleaseVersion)
	if err == nil {
		s.semver = semver
	}

	s.semverParsed = true

	return s.semver
}
