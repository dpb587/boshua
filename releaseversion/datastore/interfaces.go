package datastore

import (
	"github.com/dpb587/boshua/releaseversion"
)

type Factory interface {
	Create(provider, name string, options map[string]interface{}) (Index, error)
}

type Index interface {
	Filter(f *FilterParams) ([]releaseversion.Artifact, error)
}
