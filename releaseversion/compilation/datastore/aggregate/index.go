package aggregate

import (
	"fmt"

	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
)

type Index struct {
	indices []datastore.Index
}

var _ datastore.Index = &Index{}

func New(indices ...datastore.Index) *Index {
	return &Index{
		indices: indices,
	}
}

func (i *Index) GetCompilationArtifacts(f datastore.FilterParams) ([]compilation.Artifact, error) {
	// TODO merging behavior?
	var results []compilation.Artifact

	for idxIdx, idx := range i.indices {
		found, err := idx.GetCompilationArtifacts(f)
		if err != nil {
			return nil, fmt.Errorf("filtering %d: %v", idxIdx, err)
		}

		results = append(results, found...)
	}

	return results, nil
}

func (i *Index) StoreCompilationArtifact(artifact compilation.Artifact) error {
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
