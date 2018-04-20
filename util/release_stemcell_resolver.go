package util

import (
	"github.com/dpb587/boshua/releaseversion"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/stemcellversion"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
)

type ReleaseStemcellResolver struct {
	releaseVersionIndex  releaseversiondatastore.Index
	stemcellVersionIndex stemcellversiondatastore.Index
}

func NewReleaseStemcellResolver(
	releaseVersionIndex releaseversiondatastore.Index,
	stemcellVersionIndex stemcellversiondatastore.Index,
) *ReleaseStemcellResolver {
	return &ReleaseStemcellResolver{
		releaseVersionIndex:  releaseVersionIndex,
		stemcellVersionIndex: stemcellVersionIndex,
	}
}

func (rsr *ReleaseStemcellResolver) Resolve(releaseRef releaseversion.Reference, stemcellRef stemcellversion.Reference) (releaseversion.Subject, stemcellversion.Subject, error) {
	release, err := rsr.releaseVersionIndex.Find(releaseRef)
	if err != nil {
		return releaseversion.Subject{}, stemcellversion.Subject{}, err
	}

	stemcell, err := rsr.stemcellVersionIndex.Find(stemcellRef)
	if err != nil {
		return releaseversion.Subject{}, stemcellversion.Subject{}, err
	}

	return release, stemcell, nil
}
