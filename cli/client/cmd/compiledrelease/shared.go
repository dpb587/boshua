package compiledrelease

import (
	"fmt"

	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	"github.com/dpb587/bosh-compiled-releases/util"
	"github.com/dpb587/metalink"
)

func createMetalink(resInfo *models.CRVInfoResponse) metalink.Metalink {
	meta4 := metalink.Metalink{
		Files: []metalink.File{
			{
				Name:    fmt.Sprintf("%s-%s-on-%s-version-%s.tgz", resInfo.Data.Release.Name, resInfo.Data.Release.Version, resInfo.Data.Stemcell.OS, resInfo.Data.Stemcell.Version),
				Version: resInfo.Data.Release.Version,
				URLs: []metalink.URL{
					{
						URL: resInfo.Data.Tarball.URL,
					},
				},
			},
		},
		Generator: "bosh-compiled-releases/0.0.0",
	}

	if resInfo.Data.Tarball.Size != nil {
		meta4.Files[0].Size = *resInfo.Data.Tarball.Size
	}

	if resInfo.Data.Tarball.Published != nil {
		meta4.Published = resInfo.Data.Tarball.Published
	}

	for _, checksum := range resInfo.Data.Tarball.Checksums {
		csType, err := util.MetalinkHashType(checksum.Algorithm())
		if err != nil {
			continue
		}

		meta4.Files[0].Hashes = append(meta4.Files[0].Hashes, metalink.Hash{
			Type: csType,
			Hash: checksum.Data(),
		})
	}

	return meta4
}
