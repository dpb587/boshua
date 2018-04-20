package compiledreleaseversions

import (
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore"
)

type CompiledReleaseVersionRef struct {
	Release  releaseversions.ReleaseVersionRef
	Stemcell stemcellversions.StemcellVersionRef
}
