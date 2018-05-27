package aggregate

import (
	"fmt"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/pkg/errors"
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

func (i *index) GetAnalysisDatastore(ref stemcellversion.Reference) (analysisdatastore.Index, error) {
	for idxIdx, idx := range i.aggregated {
		analysisIndex, err := idx.GetAnalysisDatastore(ref)
		if err == analysisdatastore.UnsupportedOperationErr {
			continue
		} else if err != nil {
			return nil, errors.Wrapf(err, "getting analysis datastore", idxIdx)
		}

		return analysisIndex, nil
	}

	return nil, datastore.UnsupportedOperationErr
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
