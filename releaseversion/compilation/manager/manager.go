package manager

import (
	"github.com/dpb587/boshua/releaseversion/compilation"
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

func (rsr *Manager) Resolve(ref compilation.Reference) (releaseversion.Artifact, osversion.Artifact, error) {
	releases, err := rsr.releaseVersionIndex.Filter(ref.ReleaseVersion)
	if err != nil {
		return releaseversion.Artifact{}, osversion.Artifact{}, err
	} else if len(releases) == 0 {
		return releaseversion.Artifact{}, osversion.Artifact{}, releaseversiondatastore.NoMatchErr
	} else if len(releases) > 1 {
		return releaseversion.Artifact{}, osversion.Artifact{}, releaseversiondatastore.MultipleMatchErr
	}

	release := releases[0]

	oses, err := rsr.osVersionIndex.Filter(ref.OSVersion)
	if err != nil {
		return releaseversion.Artifact{}, osversion.Artifact{}, err
	} else if len(oses) == 0 {
		return releaseversion.Artifact{}, osversion.Artifact{}, osversiondatastore.NoMatchErr
	} else if len(oses) > 1 {
		return releaseversion.Artifact{}, osversion.Artifact{}, osversiondatastore.MultipleMatchErr
	}

	os := oses[0]

	return release, os, nil
}
