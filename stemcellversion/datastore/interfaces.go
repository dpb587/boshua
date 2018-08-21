package datastore

import (
	"github.com/dpb587/boshua/stemcellversion"
)

type ProviderName string

type Factory interface {
	Create(provider ProviderName, name string, options map[string]interface{}) (Index, error)
}

type NamedGetter func(name string) (Index, error)

type Index interface {
	GetName() string
	GetArtifacts(f FilterParams) ([]stemcellversion.Artifact, error)
	FlushCache() error
}
