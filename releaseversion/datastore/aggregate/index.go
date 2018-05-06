package aggregate

import (
	"fmt"

	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
)

type index struct {
	aggregated []datastore.Index
}

var _ datastore.Index = &index{}

func New(aggregated ...datastore.Index) datastore.Index {
	return &index{
		aggregated: aggregated,
	}
}

func (i *index) Filter(ref releaseversion.Reference) ([]releaseversion.Artifact, error) {
	var result []releaseversion.Artifact

	for idxIdx, idx := range i.aggregated {
		listed, err := idx.Filter(ref)
		if err != nil {
			return nil, fmt.Errorf("filtering %d: %v", idxIdx, err)
		}

		result = append(result, listed...)
	}

	return result, nil
}

func (i *index) Find(ref releaseversion.Reference) (releaseversion.Artifact, error) {
	return datastore.FilterForOne(i, ref)
}
