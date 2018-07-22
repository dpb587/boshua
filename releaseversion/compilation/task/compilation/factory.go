package compilation

import (
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation"
)

func New(releaseVersion releaseversion.Artifact, osVersion osversion.Artifact) *Task {
	artifact := compilation.Artifact{
		Reference: compilation.Reference{
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
