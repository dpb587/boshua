package aggregate

import (
	"fmt"

	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
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

func (i *Index) Filter(f *datastore.FilterParams) ([]stemcellversion.Artifact, error) {
	// TODO merging behavior
	var results []stemcellversion.Artifact

	for idxIdx, idx := range i.indices {
		found, err := idx.Filter(f)
		if err != nil {
			return nil, fmt.Errorf("filtering %d: %v", idxIdx, err)
		}

		results = append(results, found...)
	}

	return results, nil
}
