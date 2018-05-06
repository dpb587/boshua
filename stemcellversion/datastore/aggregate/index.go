package aggregate

import (
	"fmt"

	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
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

func (i *index) Filter(ref stemcellversion.Reference) ([]stemcellversion.Artifact, error) {
	var results []stemcellversion.Artifact

	for idxIdx, idx := range i.aggregated {
		found, err := idx.Filter(ref)
		if err != nil {
			return nil, fmt.Errorf("filtering %d: %v", idxIdx, err)
		}

		results = append(results, found...)
	}

	return results, nil
}

func (i *index) Find(ref stemcellversion.Reference) (stemcellversion.Artifact, error) {
	return datastore.FilterForOne(i, ref)
}

func (i *index) List() ([]stemcellversion.Artifact, error) {
	var results []stemcellversion.Artifact

	for idxIdx, idx := range i.aggregated {
		found, err := idx.List()
		if err != nil {
			return nil, fmt.Errorf("filtering %d: %v", idxIdx, err)
		}

		results = append(results, found...)
	}

	return results, nil
}
