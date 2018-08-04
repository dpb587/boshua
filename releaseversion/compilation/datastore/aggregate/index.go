package aggregate

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	"github.com/dpb587/metalink"
)

type index struct {
	indices []datastore.Index
}

var _ datastore.Index = &index{}

func New(indices ...datastore.Index) datastore.AnalysisIndex {
	return &index{
		indices: indices,
	}
}

func (i *index) GetCompilationArtifacts(f datastore.FilterParams) ([]compilation.Artifact, error) {
	// TODO merging behavior?
	var results []compilation.Artifact

	for idxIdx, idx := range i.indices {
		found, err := idx.GetCompilationArtifacts(f)
		if err == datastore.UnsupportedOperationErr {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("filtering %d: %v", idxIdx, err)
		}

		results = append(results, found...)
	}

	return results, nil
}

func (i *index) StoreCompilationArtifact(artifact compilation.Artifact) error {
	for idxIdx, idx := range i.indices {
		err := idx.StoreCompilationArtifact(artifact)
		if err == datastore.UnsupportedOperationErr {
			continue
		} else if err != nil {
			return fmt.Errorf("storing %d: %v", idxIdx, err)
		}

		return nil
	}

	return datastore.UnsupportedOperationErr
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
		return nil, analysisdatastore.UnsupportedOperationErr
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

	return analysisdatastore.UnsupportedOperationErr
}

func (i *index) FlushCompilationCache() error {
	for idxIdx, idx := range i.indices {
		err := idx.FlushCompilationCache()
		if err != nil {
			return fmt.Errorf("flushing %d: %v", idxIdx, err)
		}
	}

	return nil
}

func (i *index) FlushAnalysisCache() error {
	for idxIdx, idx := range i.indices {
		analysisIdx, analysisSupported := idx.(analysisdatastore.Index)
		if !analysisSupported {
			continue
		}

		err := analysisIdx.FlushAnalysisCache()
		if err != nil {
			return fmt.Errorf("flushing %d: %v", idxIdx, err)
		}
	}

	return nil
}
