package datastore

import (
	"io"

	"github.com/dpb587/boshua/analysis"
)

type Index interface {
	Filter(analysis.Reference) ([]analysis.Artifact, error)
	Store(analysis.Analyzer, analysis.Subject, io.Reader) error
}
