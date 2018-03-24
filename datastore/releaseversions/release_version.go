package releaseversions

type ReleaseVersion struct {
	ReleaseVersionRef

	Checksums      Checksums
	MetalinkSource map[string]interface{}
}
