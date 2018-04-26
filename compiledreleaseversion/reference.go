package compiledreleaseversion

import (
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/osversion"
)

type Reference struct {
	ReleaseVersion  releaseversion.Reference
	OSVersion osversion.Reference
}
