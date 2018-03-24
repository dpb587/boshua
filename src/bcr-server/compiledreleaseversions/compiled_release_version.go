package compiledreleaseversions

import "bcr-server/releaseversions"

type CompiledReleaseVersion struct {
	CompiledReleaseVersionRef

	TarballURL       string
	TarballChecksums releaseversions.Checksums
}
