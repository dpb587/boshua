package datastore

import (
	"github.com/dpb587/boshua/releaseversion"
)

type Index interface {
	Filter(releaseversion.Reference) ([]releaseversion.Artifact, error)
}
