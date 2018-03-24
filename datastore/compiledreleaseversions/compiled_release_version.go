package compiledreleaseversions

import "github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"

type CompiledReleaseVersion struct {
	CompiledReleaseVersionRef

	TarballURL       string
	TarballChecksums releaseversions.Checksums
}
