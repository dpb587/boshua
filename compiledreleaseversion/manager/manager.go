package manager

import (
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/releaseversion"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/stemcellversion"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
)

type Manager struct {
	releaseVersionIndex  releaseversiondatastore.Index
	stemcellVersionIndex stemcellversiondatastore.Index
}

func NewManager(
	releaseVersionIndex releaseversiondatastore.Index,
	stemcellVersionIndex stemcellversiondatastore.Index,
) *Manager {
	return &Manager{
		releaseVersionIndex:  releaseVersionIndex,
		stemcellVersionIndex: stemcellVersionIndex,
	}
}

func (rsr *Manager) Resolve(ref compiledreleaseversion.Reference) (releaseversion.Artifact, stemcellversion.Artifact, error) {
	release, err := rsr.releaseVersionIndex.Find(ref.ReleaseVersion)
	if err != nil {
		return releaseversion.Artifact{}, stemcellversion.Artifact{}, err
	}

	stemcell, err := rsr.stemcellVersionIndex.Find(ref.StemcellVersion)
	if err != nil {
		return releaseversion.Artifact{}, stemcellversion.Artifact{}, err
	}

	return release, stemcell, nil
}
