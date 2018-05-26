package datastore

import (
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion"
)

type Index interface {
	Find(compiledreleaseversion.Reference) (compiledreleaseversion.Artifact, error)
	Filter(compiledrelease compiledreleaseversion.Reference) ([]compiledreleaseversion.Artifact, error)
	Store(compiledreleaseversion.Artifact) error
	GetAnalysisDatastore(compiledreleaseversion.Reference) (datastore.Index, error)
}
