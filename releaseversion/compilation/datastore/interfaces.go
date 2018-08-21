package datastore

import (
	"github.com/dpb587/boshua/releaseversion/compilation"
)

type ProviderName string

type Factory interface {
	Create(provider ProviderName, name string, options map[string]interface{}) (Index, error)
}

type NamedGetter func(name string) (Index, error)

type Index interface {
	GetName() string
	GetCompilationArtifacts(f FilterParams) ([]compilation.Artifact, error)
	StoreCompilationArtifact(compilation.Artifact) error
	FlushCompilationCache() error
}
