package compilation

import (
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/stemcellversion"
)

func New(releaseVersion releaseversion.Artifact, stemcellVersion stemcellversion.Artifact) *Task {
	artifact := compiledreleaseversion.Artifact{
		ReleaseVersion:  releaseVersion.Reference,
		StemcellVersion: stemcellVersion.Reference,
	}

	return &Task{
		artifact:        artifact,
		releaseVersion:  releaseVersion,
		stemcellVersion: stemcellVersion,
	}
}
