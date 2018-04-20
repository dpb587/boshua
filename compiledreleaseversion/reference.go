package compiledreleaseversion

import (
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/stemcellversion"
)

type Reference struct {
	Release  releaseversion.Reference
	Stemcell stemcellversion.Reference
}
