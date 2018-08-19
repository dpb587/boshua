package aggregate

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
)

type index struct {
	name    string
	indices []datastore.Index
}

var _ datastore.Index = &index{}

func New(name string, indices ...datastore.Index) datastore.AnalysisIndex {
	return &index{
		name:    name,
		indices: indices,
	}
}

func (i *index) GetArtifacts(f datastore.FilterParams) ([]releaseversion.Artifact, error) {
	aggregateResults := map[string][]releaseversion.Artifact{}

	for indexIdx, index := range i.indices {
		results, err := index.GetArtifacts(f)
		if err != nil {
			return nil, errors.Wrapf(err, "filtering %d", indexIdx)
		}

		for _, result := range results {
			key := fmt.Sprintf("%s/%s", result.Name, result.Version)

			aggregateResults[key] = append(aggregateResults[key], result)
		}
	}

	var results []releaseversion.Artifact

	for aggregateResultIdx, aggregateResult := range aggregateResults {
		if len(aggregateResult) == 1 {
			results = append(results, aggregateResult...)

			continue
		}

		aggregatedResult, err := i.merge(aggregateResult)
		if err != nil {
			fmt.Printf("%#+v\n", aggregateResult)
			return nil, errors.Wrapf(err, "failed merging results for '%s'", aggregateResultIdx)
		}

		results = append(results, aggregatedResult)
	}

	for resultIdx := range results {
		results[resultIdx].Datastore = i.name
	}

	return results, nil
}

func (i *index) merge(results []releaseversion.Artifact) (releaseversion.Artifact, error) {
	// assume Name and Version already match
	result := results[0]

	for _, subresult := range results[1:] {
		if len(result.SourceTarball.Hashes) > 0 && len(subresult.SourceTarball.Hashes) > 0 {
			// TODO make this smarter
			// TODO configurable error handling; e.g. ignore vs error
			return result, nil
			// return releaseversion.Artifact{}, errors.New("multiple results with hashes found")
		}

		for _, hash := range subresult.SourceTarball.Hashes {
			// TODO avoid duplicates
			result.SourceTarball.Hashes = append(result.SourceTarball.Hashes, hash)
		}

		for _, url := range subresult.SourceTarball.URLs {
			// TODO avoid duplicates
			result.SourceTarball.URLs = append(result.SourceTarball.URLs, url)
		}

		for _, metaurl := range subresult.SourceTarball.MetaURLs {
			// TODO avoid duplicates
			result.SourceTarball.MetaURLs = append(result.SourceTarball.MetaURLs, metaurl)
		}

		// TODO handle other metalink fields
	}

	return result, nil
}

func (i *index) GetLabels() ([]string, error) {
	labelsMap := map[string]struct{}{}

	for indexIdx, idx := range i.indices {
		labels, err := idx.GetLabels()
		if err != nil {
			return nil, errors.Wrapf(err, "getting labels for %d", indexIdx)
		}

		for _, label := range labels {
			labelsMap[label] = struct{}{}
		}
	}

	var labels []string

	for label := range labelsMap {
		labels = append(labels, label)
	}

	return labels, nil
}

func (i *index) GetAnalysisArtifacts(ref analysis.Reference) ([]analysis.Artifact, error) {
	var results []analysis.Artifact
	var supported bool

	for idxIdx, idx := range i.indices {
		analysisIdx, analysisSupported := idx.(analysisdatastore.Index)
		if !analysisSupported {
			continue
		}

		supported = true

		found, err := analysisIdx.GetAnalysisArtifacts(ref)
		if err != nil {
			return nil, fmt.Errorf("analysis %d: %v", idxIdx, err)
		}

		if len(found) > 0 {
			// TODO merging behavior instead?
			return found, nil
		}
	}

	if !supported {
		return nil, analysisdatastore.UnsupportedOperationErr
	}

	return results, nil
}

func (i *index) StoreAnalysisResult(ref analysis.Reference, meta4 metalink.Metalink) error {
	for idxIdx, idx := range i.indices {
		analysisIdx, analysisSupported := idx.(analysisdatastore.Index)
		if !analysisSupported {
			continue
		}

		err := analysisIdx.StoreAnalysisResult(ref, meta4)
		if err != nil {
			return fmt.Errorf("storing %d: %v", idxIdx, err)
		}

		return nil
	}

	return analysisdatastore.UnsupportedOperationErr
}

func (i *index) FlushCache() error {
	for idxIdx, idx := range i.indices {
		err := idx.FlushCache()
		if err != nil {
			return fmt.Errorf("flushing %d: %v", idxIdx, err)
		}
	}

	return nil
}

func (i *index) FlushAnalysisCache() error {
	for idxIdx, idx := range i.indices {
		analysisIdx, analysisSupported := idx.(analysisdatastore.Index)
		if !analysisSupported {
			continue
		}

		err := analysisIdx.FlushAnalysisCache()
		if err != nil {
			return fmt.Errorf("flushing %d: %v", idxIdx, err)
		}
	}

	return nil
}
