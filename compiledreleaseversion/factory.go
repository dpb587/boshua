package compiledreleaseversion

import (
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/metalink"
)

func New(releaseVersion releaseversion.Reference, osVersion osversion.Reference, meta4File metalink.File, meta4Source map[string]interface{}) Artifact {
	return Artifact{
		ReleaseVersion:  releaseVersion,
		OSVersion: osVersion,
		MetalinkFile:    meta4File,
		MetalinkSource:  meta4Source,
	}
}
