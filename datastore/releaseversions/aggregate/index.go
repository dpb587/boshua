package aggregate

import (
	"fmt"

	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
)

type index struct {
	aggregated []releaseversions.Index
}

func New(aggregated ...releaseversions.Index) releaseversions.Index {
	return &index{
		aggregated: aggregated,
	}
}

func (i *index) List() ([]releaseversions.ReleaseVersion, error) {
	var result []releaseversions.ReleaseVersion

	for idxIdx, idx := range i.aggregated {
		listed, err := idx.List()
		if err != nil {
			return nil, fmt.Errorf("listing %d: %v", idxIdx, err)
		}

		result = append(result, listed...)
	}

	return result, nil
}

func (i *index) Find(ref releaseversions.ReleaseVersionRef) (releaseversions.ReleaseVersion, error) {
	for idxIdx, idx := range i.aggregated {
		found, err := idx.Find(ref)
		if err == releaseversions.MissingErr {
			continue
		} else if err != nil {
			return releaseversions.ReleaseVersion{}, fmt.Errorf("listing %d: %v", idxIdx, err)
		}

		return found, nil
	}

	return releaseversions.ReleaseVersion{}, releaseversions.MissingErr
}