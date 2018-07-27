package aggregate

import (
	"fmt"

	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/pkg/errors"
)

type Index struct {
	indices []datastore.Index
}

var _ datastore.Index = &Index{}

func New(indices ...datastore.Index) *Index {
	return &Index{
		indices: indices,
	}
}

func (i *Index) Filter(f *datastore.FilterParams) ([]releaseversion.Artifact, error) {
	aggregateResults := map[string][]releaseversion.Artifact{}

	for indexIdx, index := range i.indices {
		results, err := index.Filter(f)
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

	return results, nil
}

func (i *Index) merge(results []releaseversion.Artifact) (releaseversion.Artifact, error) {
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

func (i *Index) Labels() ([]string, error) {
	labelsMap := map[string]struct{}{}

	for indexIdx, idx := range i.indices {
		labels, err := idx.Labels()
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
