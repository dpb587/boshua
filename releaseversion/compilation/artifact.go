package compilation

import (
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/metalink"
)

type Artifact struct {
	Datastore string `json:"-"`

	Release releaseversion.Reference `json:"release"`
	OS      osversion.Reference      `json:"os"`

	Labels []string `json:"labels"`

	Tarball metalink.File `json:"tarball"`
}

func (s Artifact) MetalinkFile() metalink.File {
	return s.Tarball
}

func (s Artifact) Reference() interface{} {
	return Reference{
		ReleaseVersion: s.Release,
		OSVersion:      s.OS,
	}
}

func (s Artifact) GetLabels() []string {
	return s.Labels
}

func (s Artifact) GetDatastoreName() string {
	return s.Datastore
}
