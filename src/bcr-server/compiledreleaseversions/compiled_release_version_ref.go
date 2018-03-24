package compiledreleaseversions

import (
	"bcr-server/releaseversions"
	"bcr-server/stemcellversions"
)

type CompiledReleaseVersionRef struct {
	Release  releaseversions.ReleaseVersionRef
	Stemcell stemcellversions.StemcellVersionRef
}
