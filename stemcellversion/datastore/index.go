package datastore

import "github.com/dpb587/boshua/stemcellversion"

type Index interface {
	Find(ref stemcellversion.Reference) (stemcellversion.Subject, error)
	List() ([]stemcellversion.Subject, error)
}
