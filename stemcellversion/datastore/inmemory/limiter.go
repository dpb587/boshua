package inmemory

import (
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
)

func LimitArtifacts(artifacts []stemcellversion.Artifact, l datastore.LimitParams) ([]stemcellversion.Artifact, error) {
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
	stemcellversion.Sort(artifacts)

	if l.OffsetExpected {
		artifacts = artifacts[l.Offset:]
	}

	if l.LimitExpected {
		artifacts = artifacts[:l.Limit]
	}

	return artifacts, nil
}
