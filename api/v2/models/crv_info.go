package models

type CRVInfoRequest struct {
	Data CRVInfoRequestData `json:"data"`
}

type CRVInfoRequestData struct {
	Release  ReleaseRef  `json:"release"`
	Stemcell StemcellRef `json:"stemcell"`
}

type CRVInfoResponse struct {
	Data CRVInfoResponseCompiledRelease `json:"data"`
}

type CRVInfoResponseCompiledRelease struct {
	URL       string      `json:"url"`
	Checksums []Checksum  `json:"checksums"`
	Release   ReleaseRef  `json:"release"`
	Stemcell  StemcellRef `json:"stemcell"`
}
