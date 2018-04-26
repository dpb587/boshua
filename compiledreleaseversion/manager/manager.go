package manager

import (
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/osversion"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	"github.com/dpb587/boshua/releaseversion"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
)

type Manager struct {
	releaseVersionIndex releaseversiondatastore.Index
	osVersionIndex      osversiondatastore.Index
}

func NewManager(
	releaseVersionIndex releaseversiondatastore.Index,
	osVersionIndex osversiondatastore.Index,
) *Manager {
	return &Manager{
		releaseVersionIndex: releaseVersionIndex,
		osVersionIndex:      osVersionIndex,
	}
}

func (rsr *Manager) Resolve(ref compiledreleaseversion.Reference) (releaseversion.Artifact, osversion.Artifact, error) {
	release, err := rsr.releaseVersionIndex.Find(ref.ReleaseVersion)
	if err != nil {
		return releaseversion.Artifact{}, osversion.Artifact{}, err
	}

	os, err := rsr.osVersionIndex.Find(ref.OSVersion)
	if err != nil {
		return releaseversion.Artifact{}, osversion.Artifact{}, err
	}

	return release, os, nil
}
