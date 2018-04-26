package compiledreleaseversion

import (
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/metalink"
)

func New(releaseVersion releaseversion.Reference, stemcellVersion stemcellversion.Reference, meta4File metalink.File, meta4Source map[string]interface{}) Artifact {
	return Artifact{
		ReleaseVersion:  releaseVersion,
		StemcellVersion: stemcellVersion,
		MetalinkFile:    meta4File,
		MetalinkSource:  meta4Source,
	}
}
