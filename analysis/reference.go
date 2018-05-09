package analysis

import (
	"crypto/sha1"
	"fmt"

	"github.com/dpb587/boshua/artifact"
)

type Reference struct {
	Artifact artifact.ArtifactReference
	Analyzer string
}

var _ artifact.ArtifactReference = &Reference{}

func (r Reference) ArtifactReference() artifact.Reference {
	return artifact.Reference{
		Context: "analysis",
		ID:      r.id(),
	}
}

func (r Reference) ArtifactStorageDir() string {
	return fmt.Sprintf("%s/analysis/%s", r.Artifact.ArtifactStorageDir(), r.Analyzer)
}

func (r Reference) id() string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf(
		"analysis:v1:%s:%s",
		r.Artifact.ArtifactReference().String(),
		r.Analyzer,
	)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
