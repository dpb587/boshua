package compilation

import (
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/osversion"
)

func New(releaseVersion releaseversion.Artifact, osVersion osversion.Artifact) *Task {
	artifact := compiledreleaseversion.Artifact{
		ReleaseVersion:  releaseVersion.Reference,
		OSVersion: osVersion.Reference,
	}

	return &Task{
		artifact:        artifact,
		releaseVersion:  releaseVersion,
		osVersion: osVersion,
	}
}
