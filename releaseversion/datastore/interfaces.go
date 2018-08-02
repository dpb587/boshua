package datastore

import (
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/releaseversion"
)

type Factory interface {
	Create(provider, name string, options map[string]interface{}) (Index, error)
}

type Index interface {
	GetArtifacts(f FilterParams) ([]releaseversion.Artifact, error)
	GetLabels() ([]string, error)
	FlushCache() error
}

type AnalysisIndex interface {
	Index
	analysisdatastore.Index
}
