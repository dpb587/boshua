package datastore

import "github.com/dpb587/boshua/osversion"

type Index interface {
	Filter(osversion.Reference) ([]osversion.Artifact, error)
}
