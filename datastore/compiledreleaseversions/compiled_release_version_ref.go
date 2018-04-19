package compiledreleaseversions

import (
	"github.com/dpb587/boshua/datastore/releaseversions"
	"github.com/dpb587/boshua/datastore/stemcellversions"
)

type CompiledReleaseVersionRef struct {
	Release  releaseversions.ReleaseVersionRef
	Stemcell stemcellversions.StemcellVersionRef
}
