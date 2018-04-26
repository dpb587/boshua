package compiledreleaseversion

import (
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/stemcellversion"
)

type Reference struct {
	ReleaseVersion  releaseversion.Reference
	StemcellVersion stemcellversion.Reference
}
