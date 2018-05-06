package aggregate

import (
	"fmt"

	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
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

func (i *index) Filter(ref compiledreleaseversion.Reference) ([]compiledreleaseversion.Artifact, error) {
	var results []compiledreleaseversion.Artifact

	for idxIdx, idx := range i.aggregated {
		found, err := idx.Filter(ref)
		if err != nil {
			return nil, fmt.Errorf("listing %d: %v", idxIdx, err)
		}

		results = append(results, found...)
	}

	return results, nil
}

func (i *index) Find(ref compiledreleaseversion.Reference) (compiledreleaseversion.Artifact, error) {
	return datastore.FilterForOne(i, ref)
}
