package datastore

import (
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/releaseversion/compilation"
)

type ProviderName string

type Factory interface {
	Create(provider ProviderName, name string, options map[string]interface{}) (Index, error)
}

type NamedGetter func(name string) (Index, error)

type Index interface {
	GetCompilationArtifacts(f FilterParams) ([]compilation.Artifact, error)
	StoreCompilationArtifact(compilation.Artifact) error
	FlushCompilationCache() error
}

type AnalysisIndex interface {
	Index
	analysisdatastore.Index
}
