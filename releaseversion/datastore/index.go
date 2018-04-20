package datastore

import "github.com/dpb587/boshua/releaseversion"

type Index interface {
	Find(ref releaseversion.Reference) (releaseversion.Subject, error)
	List() ([]releaseversion.Subject, error)
}
