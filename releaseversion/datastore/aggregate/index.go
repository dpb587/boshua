package aggregate

import (
	"fmt"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/analysis/datastore/localcache"
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
	var results []releaseversion.Artifact

	for idxIdx, idx := range i.aggregated {
		found, err := idx.Filter(ref)
		if err != nil {
			return nil, fmt.Errorf("filtering %d: %v", idxIdx, err)
		}

		results = append(results, found...)
	}

	return results, nil
}

func (i *index) Find(ref releaseversion.Reference) (releaseversion.Artifact, error) {
	return datastore.FilterForOne(i, ref)
}

func (i *index) GetAnalysisDatastore() analysisdatastore.Index { // TODO aggregate probably requires err for Unsupported check
	return localcache.New()
}
