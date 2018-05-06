package datastore

import "github.com/dpb587/boshua/analysis"

type Index interface {
	Filter(analysis.Reference) ([]analysis.Artifact, error)
}
