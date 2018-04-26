package boshua

import (
	"fmt"

	"github.com/dpb587/metalink"
)

type Artifact interface {
	ArtifactReference() Reference
	ArtifactMetalink() metalink.Metalink
	ArtifactMetalinkStorage() map[string]interface{}
}

type Reference struct {
	Context string
	ID      string
}

func (r Reference) String() string {
	return fmt.Sprintf("%s/%s", r.Context, r.ID)
}
