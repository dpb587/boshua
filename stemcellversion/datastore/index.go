package datastore

import (
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/stemcellversion"
)

type Index interface {
	Find(stemcellversion.Reference) (stemcellversion.Artifact, error)
	Filter(stemcellversion.Reference) ([]stemcellversion.Artifact, error)
	GetAnalysisDatastore(stemcellversion.Reference) (datastore.Index, error)

	// TODO remove; kept for osversion
	List() ([]stemcellversion.Artifact, error)
}
