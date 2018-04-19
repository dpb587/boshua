package aggregate

import (
	"fmt"

	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions"
)

type index struct {
	aggregated []compiledreleaseversions.Index
}

var _ compiledreleaseversions.Index = &index{}

func New(aggregated ...compiledreleaseversions.Index) compiledreleaseversions.Index {
	return &index{
		aggregated: aggregated,
	}
}

func (i *index) List() ([]compiledreleaseversions.CompiledReleaseVersion, error) {
	var result []compiledreleaseversions.CompiledReleaseVersion

	for idxIdx, idx := range i.aggregated {
		listed, err := idx.List()
		if err != nil {
			return nil, fmt.Errorf("listing %d: %v", idxIdx, err)
		}

		result = append(result, listed...)
	}

	return result, nil
}

func (i *index) Find(ref compiledreleaseversions.CompiledReleaseVersionRef) (compiledreleaseversions.CompiledReleaseVersion, error) {
	for idxIdx, idx := range i.aggregated {
		found, err := idx.Find(ref)
		if err == compiledreleaseversions.MissingErr {
			continue
		} else if err != nil {
			return compiledreleaseversions.CompiledReleaseVersion{}, fmt.Errorf("listing %d: %v", idxIdx, err)
		}

		return found, nil
	}

	return compiledreleaseversions.CompiledReleaseVersion{}, compiledreleaseversions.MissingErr
}
