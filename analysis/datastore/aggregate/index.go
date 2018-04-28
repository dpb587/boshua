package aggregate

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
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

func (i *index) List() ([]analysis.Artifact, error) {
	var result []analysis.Artifact

	for idxIdx, idx := range i.aggregated {
		listed, err := idx.List()
		if err != nil {
			return nil, fmt.Errorf("listing %d: %v", idxIdx, err)
		}

		result = append(result, listed...)
	}

	return result, nil
}

func (i *index) Find(ref analysis.Reference) (analysis.Artifact, error) {
	for idxIdx, idx := range i.aggregated {
		found, err := idx.Find(ref)
		if err == datastore.MissingErr {
			continue
		} else if err != nil {
			return analysis.Artifact{}, fmt.Errorf("listing %d: %v", idxIdx, err)
		}

		return found, nil
	}

	return analysis.Artifact{}, datastore.MissingErr
}
