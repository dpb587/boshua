package compiledreleaseversion

import (
	"crypto/sha1"
	"fmt"

	"github.com/dpb587/boshua"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/metalink"
)

type Artifact struct {
	MetalinkFile   metalink.File
	MetalinkSource map[string]interface{}

	ReleaseVersion releaseversion.Reference
	OSVersion      osversion.Reference
}

func (s Artifact) ArtifactReference() boshua.Reference {
	return boshua.Reference{
		Context: "compiledreleaseversion",
		ID:      s.id(),
	}
}

func (s Artifact) ArtifactMetalinkStorage() map[string]interface{} {
	return map[string]interface{}{
		"uri": fmt.Sprintf("git@github.com:dpb587/bosh-compiled-releases-index.git//%s", s.StoragePath()),
		"options": map[string]string{
			"private_key": "((index_private_key))",
		},
	}
}

func (s Artifact) StoragePath() string {
	return fmt.Sprintf("%s/%s", s.ArtifactReference().Context, s.ArtifactReference().ID)
}

func (s Artifact) id() string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf(
		"compiledreleaseversion:v1:%s:%s:%s",
		s.OSVersion.Name,
		s.OSVersion.Version,
		s.ReleaseVersion.ArtifactReference().ID,
	)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
