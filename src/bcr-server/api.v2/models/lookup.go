package models

type LookupRequest struct {
	Data LookupRequestData `json:"data"`
}

type LookupRequestData struct {
	Release  ReleaseRef  `json:"release"`
	Stemcell StemcellRef `json:"stemcell"`
}

type LookupResponse struct {
	CompiledRelease LookupResponseCompiledRelease `json:"compiled_release"`
}

type LookupResponseCompiledRelease struct {
	URL       string      `json:"url"`
	Checksums []Checksum  `json:"checksums"`
	Release   ReleaseRef  `json:"release"`
	Stemcell  StemcellRef `json:"stemcell"`
}
