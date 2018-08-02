package datastore

import (
	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/metalink"
)

type analysisIndex struct {
	index         Index
	analysisIndex analysisdatastore.Index
}

var _ Index = &analysisIndex{}
var _ analysisdatastore.Index = &analysisIndex{}

func NewAnalysisIndex(idx Index, analysisIdx analysisdatastore.Index) *analysisIndex {
	return &analysisIndex{
		index:         idx,
		analysisIndex: analysisIdx,
	}
}

func (i *analysisIndex) GetArtifacts(f FilterParams) ([]stemcellversion.Artifact, error) {
	return i.index.GetArtifacts(f)
}

func (i *analysisIndex) GetAnalysisArtifacts(ref analysis.Reference) ([]analysis.Artifact, error) {
	return i.analysisIndex.GetAnalysisArtifacts(ref)
}

func (i *analysisIndex) StoreAnalysisResult(ref analysis.Reference, meta4 metalink.Metalink) error {
	return i.analysisIndex.StoreAnalysisResult(ref, meta4)
}

func (i *analysisIndex) FlushCache() error {
	return i.index.FlushCache()
}

func (i *analysisIndex) FlushAnalysisCache() error {
	return i.analysisIndex.FlushAnalysisCache()
}
