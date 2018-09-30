package pivnetfile

import (
	"github.com/Masterminds/semver"
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/metalink"
)

type Artifact struct {
	Datastore string `json:"-"`

	ProductSlug    string `json:"product_slug"`
	ReleaseID      int    `json:"release_id"`
	FileID         int    `json:"file_id"`

	ReleaseVersion string        `json:"release_version"`
	File           metalink.File `json:"file"`

	semver       *semver.Version
	semverParsed bool
}

var _ artifact.Artifact = &Artifact{}

func (s Artifact) MetalinkFile() metalink.File {
	return s.File
}

func (s Artifact) Reference() interface{} {
	return Reference{
		ProductSlug: s.ProductSlug,
		ReleaseID:   s.ReleaseID,
		FileID:      s.FileID,
	}
}

func (s Artifact) GetLabels() []string {
	// TODO ReadyToServe == stable? extract stability?
	return nil
}

func (s Artifact) GetDatastoreName() string {
	return s.Datastore
}

func (s Artifact) Semver() *semver.Version {
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
