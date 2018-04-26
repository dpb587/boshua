package datastore

import "github.com/dpb587/boshua/osversion"

type Index interface {
	Find(ref osversion.Reference) (osversion.Artifact, error)
	List() ([]osversion.Artifact, error)
}
