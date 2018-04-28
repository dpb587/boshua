package datastore

import "github.com/dpb587/boshua/analysis"

type Index interface {
	Find(analysis.Reference) (analysis.Artifact, error)
	List() ([]analysis.Artifact, error)
}
