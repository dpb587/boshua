package datastore

import "github.com/dpb587/boshua/stemcellversion"

type Index interface {
	Find(stemcellversion.Reference) (stemcellversion.Artifact, error)
	Filter(stemcellversion.Reference) ([]stemcellversion.Artifact, error)

	// TODO remove; kept for osversion
	List() ([]stemcellversion.Artifact, error)
}
