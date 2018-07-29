package aggregate

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/metalink"
)

type index struct {
	indices []datastore.Index
}

var _ datastore.Index = &index{}
var _ analysisdatastore.Index = &index{}

func New(indices ...datastore.Index) *index {
	return &index{
		indices: indices,
	}
}

func (i *index) GetArtifacts(f datastore.FilterParams) ([]stemcellversion.Artifact, error) {
	// TODO merging behavior
	var results []stemcellversion.Artifact

	for idxIdx, idx := range i.indices {
		found, err := idx.GetArtifacts(f)
		if err != nil {
			return nil, fmt.Errorf("filtering %d: %v", idxIdx, err)
		}

		results = append(results, found...)
	}

	return results, nil
}

func (i *index) GetAnalysisArtifacts(ref analysis.Reference) ([]analysis.Artifact, error) {
	var results []analysis.Artifact
	var supported bool

	for idxIdx, idx := range i.indices {
		analysisIdx, analysisSupported := idx.(analysisdatastore.Index)
		if !analysisSupported {
			continue
		}

		supported = true

		found, err := analysisIdx.GetAnalysisArtifacts(ref)
		if err != nil {
			return nil, fmt.Errorf("analysis %d: %v", idxIdx, err)
		}

		if len(found) > 0 {
			// TODO merging behavior instead?
			return found, nil
		}
	}

	if !supported {
		return nil, datastore.UnsupportedOperationErr
	}

	return results, nil
}

func (i *index) StoreAnalysisResult(ref analysis.Reference, meta4 metalink.Metalink) error {
	for idxIdx, idx := range i.indices {
		analysisIdx, analysisSupported := idx.(analysisdatastore.Index)
		if !analysisSupported {
			continue
		}

		err := analysisIdx.StoreAnalysisResult(ref, meta4)
		if err != nil {
			return fmt.Errorf("storing %d: %v", idxIdx, err)
		}

		return nil
	}

	return datastore.UnsupportedOperationErr
}

func (i *index) FlushCache() error {
	for idxIdx, idx := range i.indices {
		analysisIdx, analysisSupported := idx.(analysisdatastore.Index)
		if !analysisSupported {
			continue
		}

		err := analysisIdx.FlushCache()
		if err != nil {
			return fmt.Errorf("flushing %d: %v", idxIdx, err)
		}
	}

	return nil
}
