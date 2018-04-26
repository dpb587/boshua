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

func (i *index) List() ([]stemcellversion.Artifact, error) {
	var result []stemcellversion.Artifact

	for idxIdx, idx := range i.aggregated {
		listed, err := idx.List()
		if err != nil {
			return nil, fmt.Errorf("listing %d: %v", idxIdx, err)
		}

		result = append(result, listed...)
	}

	return result, nil
}

func (i *index) Find(ref stemcellversion.Reference) (stemcellversion.Artifact, error) {
	for idxIdx, idx := range i.aggregated {
		found, err := idx.Find(ref)
		if err == datastore.MissingErr {
			continue
		} else if err != nil {
			return stemcellversion.Artifact{}, fmt.Errorf("listing %d: %v", idxIdx, err)
		}

		return found, nil
	}

	return stemcellversion.Artifact{}, datastore.MissingErr
}
