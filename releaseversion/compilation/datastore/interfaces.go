package datastore

import (
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/releaseversion/compilation"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
)

type Factory interface {
	Create(provider, name string, options map[string]interface{}, releaseVersionIndex releaseversiondatastore.Index) (Index, error)
}

type Index interface {
	GetCompilationArtifacts(f FilterParams) ([]compilation.Artifact, error)
	StoreCompilationArtifact(compilation.Artifact) error
}

type AnalysisIndex interface {
	Index
	analysisdatastore.Index
}
