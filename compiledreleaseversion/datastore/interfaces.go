package datastore

import (
	"github.com/dpb587/boshua/compiledreleaseversion"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
)

type Factory interface {
	Create(provider, name string, options map[string]interface{}, releaseVersionIndex releaseversiondatastore.Index) (Index, error)
}

type Index interface {
	Filter(f *FilterParams) ([]compiledreleaseversion.Artifact, error)
}
