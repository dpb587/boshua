package compiledreleaseversions

import (
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
)

type CompiledReleaseVersionRef struct {
	Release  releaseversions.ReleaseVersionRef
	Stemcell stemcellversions.StemcellVersionRef
}
