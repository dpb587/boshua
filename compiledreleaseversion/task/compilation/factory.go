package compilation

import (
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
)

func New(releaseVersion releaseversion.Artifact, osVersion osversion.Artifact) *Task {
	artifact := compiledreleaseversion.Artifact{
		Reference: compiledreleaseversion.Reference{
			ReleaseVersion: releaseVersion.Reference,
			OSVersion:      osVersion.Reference,
		},
	}

	return &Task{
		artifact:       artifact,
		releaseVersion: releaseVersion,
		osVersion:      osVersion,
	}
}
