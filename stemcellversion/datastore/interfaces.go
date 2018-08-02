package datastore

import (
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/stemcellversion"
)

type Factory interface {
	Create(provider, name string, options map[string]interface{}) (Index, error)
}

type Index interface {
	GetArtifacts(f FilterParams) ([]stemcellversion.Artifact, error)
	FlushCache() error
}

type AnalysisIndex interface {
	Index
	analysisdatastore.Index
}
