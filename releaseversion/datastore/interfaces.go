package datastore

import (
	"github.com/dpb587/boshua/releaseversion"
)

type Factory interface {
	Create(provider, name string, options map[string]interface{}) (Index, error)
}

type Index interface {
	GetArtifacts(f FilterParams) ([]releaseversion.Artifact, error)
	GetLabels() ([]string, error)
}
