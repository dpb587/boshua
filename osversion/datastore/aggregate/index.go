package aggregate

import (
	"fmt"

	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/osversion/datastore"
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

func (i *index) List() ([]osversion.Artifact, error) {
	var result []osversion.Artifact

	for idxIdx, idx := range i.aggregated {
		listed, err := idx.List()
		if err != nil {
			return nil, fmt.Errorf("listing %d: %v", idxIdx, err)
		}

		result = append(result, listed...)
	}

	return result, nil
}

func (i *index) Find(ref osversion.Reference) (osversion.Artifact, error) {
	for idxIdx, idx := range i.aggregated {
		found, err := idx.Find(ref)
		if err == datastore.MissingErr {
			continue
		} else if err != nil {
			return osversion.Artifact{}, fmt.Errorf("listing %d: %v", idxIdx, err)
		}

		return found, nil
	}

	return osversion.Artifact{}, datastore.MissingErr
}
