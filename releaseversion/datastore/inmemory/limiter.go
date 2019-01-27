package inmemory

import (
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
)

func LimitArtifacts(artifacts []releaseversion.Artifact, l datastore.LimitParams) ([]releaseversion.Artifact, error) {
	if artifacts == nil {
		return artifacts, nil
	}

	artifactsLen := len(artifacts)

	// min/max BEFORE limiting; used to assert quality of results before getting them
	if l.MinExpected {
		if len(artifacts) < l.Min {
			return nil, datastore.NewUnexpectedMinCountError(l.Min, artifactsLen)
		}
	}

	if l.MaxExpected {
		if artifactsLen > l.Max {
			return nil, datastore.NewUnexpectedMaxCountError(l.Max, artifactsLen)
		}
	}

	// TODO configurable sort?
	// TODO inefficient to always sort?
	releaseversion.Sort(artifacts)

	if l.OffsetExpected {
		artifacts = artifacts[l.Offset:]
	}

	if l.LimitExpected {
		ll := l.Limit
		if lla := len(artifacts); ll > lla {
			ll = lla
		}

		artifacts = artifacts[:ll]
	}

	return artifacts, nil
}
