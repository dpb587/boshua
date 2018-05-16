package datastore

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/metalink"
)

type Index interface {
	Filter(analysis.Reference) ([]analysis.Artifact, error)
	Store(analysis.Reference, metalink.Metalink) error
}
