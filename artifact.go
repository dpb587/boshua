package boshua

import (
	"fmt"

	"github.com/dpb587/metalink"
)

type Artifact interface {
	ArtifactReference

	ArtifactMetalinkFile() metalink.File
	ArtifactMetalinkStorage() map[string]interface{}
}

type ArtifactReference interface {
	ArtifactReference() Reference
	ArtifactStorageDir() string
}

type Reference struct {
	Context string
	ID      string
}

func (r Reference) String() string {
	return fmt.Sprintf("%s/%s", r.Context, r.ID)
}
