package analysis

import (
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/metalink"
)

type Artifact struct {
	Datastore string

	reference    Reference
	metalinkFile metalink.File
}

var _ artifact.Artifact = &Artifact{}

func (s Artifact) MetalinkFile() metalink.File {
	return s.metalinkFile
}

func (s Artifact) Reference() interface{} {
	return s.reference
}

func (Artifact) GetLabels() []string {
	return nil
}

func (s Artifact) GetDatastoreName() string {
	return s.Datastore
}
