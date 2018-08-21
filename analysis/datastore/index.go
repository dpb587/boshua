package datastore

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/metalink"
)

type Index interface {
	GetName() string
	GetAnalysisArtifacts(analysis.Reference) ([]analysis.Artifact, error)
	StoreAnalysisResult(analysis.Reference, metalink.Metalink) error
	FlushAnalysisCache() error // TODO rename? intent is force reload next time
}
