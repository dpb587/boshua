package util

import (
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore"
)

type ReleaseStemcellResolver struct {
	releaseVersionIndex  releaseversions.Index
	stemcellVersionIndex stemcellversions.Index
}

func NewReleaseStemcellResolver(
	releaseVersionIndex releaseversions.Index,
	stemcellVersionIndex stemcellversions.Index,
) *ReleaseStemcellResolver {
	return &ReleaseStemcellResolver{
		releaseVersionIndex:  releaseVersionIndex,
		stemcellVersionIndex: stemcellVersionIndex,
	}
}

func (rsr *ReleaseStemcellResolver) Resolve(releaseRef releaseversions.ReleaseVersionRef, stemcellRef stemcellversions.StemcellVersionRef) (releaseversions.ReleaseVersion, stemcellversions.StemcellVersion, error) {
	release, err := rsr.releaseVersionIndex.Find(releaseRef)
	if err != nil {
		return releaseversions.ReleaseVersion{}, stemcellversions.StemcellVersion{}, err
	}

	stemcell, err := rsr.stemcellVersionIndex.Find(stemcellRef)
	if err != nil {
		return releaseversions.ReleaseVersion{}, stemcellversions.StemcellVersion{}, err
	}

	return release, stemcell, nil
}
