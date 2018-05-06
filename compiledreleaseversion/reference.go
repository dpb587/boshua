package compiledreleaseversion

import (
	"crypto/sha1"
	"fmt"

	"github.com/dpb587/boshua"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
)

type Reference struct {
	ReleaseVersion releaseversion.Reference `json:"release"`
	OSVersion      osversion.Reference      `json:"os"`
}

var _ boshua.ArtifactReference = &Reference{}

func (r Reference) ArtifactReference() boshua.Reference {
	return boshua.Reference{
		Context: "compiledreleaseversion",
		ID:      r.id(),
	}
}

func (r Reference) ArtifactStorageDir() string {
	ref := r.ArtifactReference()

	return fmt.Sprintf(
		"%s/%s/%s/%s",
		ref.Context,
		ref.ID[0:2],
		ref.ID[2:4],
		ref.ID[4:],
	)
}

func (r Reference) id() string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf(
		"compiledreleaseversion:v1:%s:%s:%s",
		r.OSVersion.Name,
		r.OSVersion.Version,
		r.ReleaseVersion.ArtifactReference().ID,
	)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
