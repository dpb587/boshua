package datastore

import (
	"github.com/dpb587/boshua/releaseversion"
)

type ProviderName string

type Factory interface {
	Create(provider ProviderName, name string, options map[string]interface{}) (Index, error)
}

type NamedGetter func(name string) (Index, error)

type Index interface {
	GetName() string
	GetArtifacts(f FilterParams, l LimitParams) ([]releaseversion.Artifact, error)
	GetLabels() ([]string, error)
	FlushCache() error
}
