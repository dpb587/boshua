package datastore

import (
	"github.com/dpb587/boshua/pivnetfile"
)

func FilterParamsFromArtifact(artifact pivnetfile.Artifact) FilterParams {
	f := FilterParams{
		ProductNameExpected: true,
		ProductName:         artifact.ProductName,

		ReleaseIDExpected: true,
		ReleaseID:         artifact.ReleaseID,

		FileIDExpected: true,
		FileID:         artifact.FileID,
	}

	return f
}
