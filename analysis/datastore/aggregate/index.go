package aggregate

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
)

type index struct {
	name       string
	aggregated []datastore.Index
}

var _ datastore.Index = &index{}

func New(name string, aggregated ...datastore.Index) datastore.Index {
	return &index{
		name:       name,
		aggregated: aggregated,
	}
}

func (i *index) GetName() string {
	return i.name
}

func (i *index) GetAnalysisArtifacts(ref analysis.Reference) ([]analysis.Artifact, error) {
	var results []analysis.Artifact

	for idxIdx, idx := range i.aggregated {
		found, err := idx.GetAnalysisArtifacts(ref)
		if err != nil {
			return nil, fmt.Errorf("listing %d: %v", idxIdx, err)
		}

		results = append(results, found...)
	}

	return results, nil
}

func (i *index) FlushAnalysisCache() error {
	for idxIdx, idx := range i.aggregated {
		err := idx.FlushAnalysisCache()
		if err != nil {
			return errors.Wrapf(err, "flushing %d", idxIdx)
		}
	}

	return nil
}

func (i *index) StoreAnalysisResult(ref analysis.Reference, artifactMeta4 metalink.Metalink) error {
	for idxIdx, idx := range i.aggregated {
		err := idx.StoreAnalysisResult(ref, artifactMeta4)
		if err == datastore.UnsupportedOperationErr {
			continue
		} else if err != nil {
			return fmt.Errorf("storing %d: %v", idxIdx, err)
		}

		return nil
	}

	return datastore.UnsupportedOperationErr
}
