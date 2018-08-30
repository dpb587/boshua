package aggregate

import (
	"fmt"

	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	"github.com/pkg/errors"
)

type index struct {
	name    string
	indices []datastore.Index
}

var _ datastore.Index = &index{}

func New(name string, indices ...datastore.Index) datastore.Index {
	return &index{
		name:    name,
		indices: indices,
	}
}

func (i *index) GetName() string {
	return i.name
}

func (i *index) GetCompilationArtifacts(f datastore.FilterParams) ([]compilation.Artifact, error) {
	// TODO merging behavior?
	var results []compilation.Artifact

	for idxIdx, idx := range i.indices {
		found, err := idx.GetCompilationArtifacts(f)
		if err == datastore.UnsupportedOperationErr {
			continue
		} else if err != nil {
			if len(i.indices) == 1 {
				return nil, err
			}

			return nil, errors.Wrapf(err, "filtering %d", idxIdx)
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

func (i *index) FlushCompilationCache() error {
	for idxIdx, idx := range i.indices {
		err := idx.FlushCompilationCache()
		if err != nil {
			return fmt.Errorf("flushing %d: %v", idxIdx, err)
		}
	}

	return nil
}
