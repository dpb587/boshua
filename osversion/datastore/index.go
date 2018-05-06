package datastore

import "github.com/dpb587/boshua/osversion"

type Index interface {
	Find(osversion.Reference) (osversion.Artifact, error)
	Filter(osversion.Reference) ([]osversion.Artifact, error)
}
