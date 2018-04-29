package datastore

import "github.com/dpb587/boshua/stemcellversion"

type Index interface {
	Find(ref stemcellversion.Reference) (stemcellversion.Artifact, error)
	List() ([]stemcellversion.Artifact, error)
}
