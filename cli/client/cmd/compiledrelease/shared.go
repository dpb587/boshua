package compiledrelease

import (
	api "github.com/dpb587/boshua/api/v2/models/compiledreleaseversion"
	"github.com/dpb587/metalink"
)

func createMetalink(compilation *api.GETCompilationResponse) metalink.Metalink {
	meta4 := metalink.Metalink{
		Files: []metalink.File{
			compilation.Data,
		},
		Generator: "bosh-compiled-releases/0.0.0",
	}

	// TODO restore
	// if resInfo.Data.Tarball.Published != nil {
	// 	meta4.Published = resInfo.Data.Tarball.Published
	// }

	return meta4
}
