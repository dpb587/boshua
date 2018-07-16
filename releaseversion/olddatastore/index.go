package datastore

import (
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/releaseversion"
)

type Index interface {
	Find(releaseversion.Reference) (releaseversion.Artifact, error)
	Filter(releaseversion.Reference) ([]releaseversion.Artifact, error)
	GetAnalysisDatastore() datastore.Index
}
