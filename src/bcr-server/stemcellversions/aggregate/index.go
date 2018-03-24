package aggregate

import (
	"bcr-server/stemcellversions"
	"fmt"
)

type index struct {
	aggregated []stemcellversions.Index
}

func New(aggregated ...stemcellversions.Index) stemcellversions.Index {
	return &index{
		aggregated: aggregated,
	}
}

func (i *index) List() ([]stemcellversions.StemcellVersion, error) {
	var result []stemcellversions.StemcellVersion

	for idxIdx, idx := range i.aggregated {
		listed, err := idx.List()
		if err != nil {
			return nil, fmt.Errorf("listing %d: %v", idxIdx, err)
		}

		result = append(result, listed...)
	}

	return result, nil
}

func (i *index) Find(ref stemcellversions.StemcellVersionRef) (stemcellversions.StemcellVersion, error) {
	for idxIdx, idx := range i.aggregated {
		found, err := idx.Find(ref)
		if err == stemcellversions.MissingErr {
			continue
		} else if err != nil {
			return stemcellversions.StemcellVersion{}, fmt.Errorf("listing %d: %v", idxIdx, err)
		}

		return found, nil
	}

	return stemcellversions.StemcellVersion{}, stemcellversions.MissingErr
}
