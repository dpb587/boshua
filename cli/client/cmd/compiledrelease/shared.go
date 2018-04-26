package compiledrelease

import (
	"github.com/dpb587/boshua/api/v2/models"
	"github.com/dpb587/metalink"
)

func createMetalink(resInfo *models.CRVInfoResponse) metalink.Metalink {
	meta4 := metalink.Metalink{
		Files: []metalink.File{
			resInfo.Data.Artifact,
		},
		Generator: "bosh-compiled-releases/0.0.0",
	}

	// TODO restore
	// if resInfo.Data.Tarball.Published != nil {
	// 	meta4.Published = resInfo.Data.Tarball.Published
	// }

	return meta4
}
