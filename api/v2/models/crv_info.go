package models

type CRVInfoStatus string

const (
	CRVInfoStatusUnknown     CRVInfoStatus = "unknown"
	CRVInfoStatusPending     CRVInfoStatus = "pending"
	CRVInfoStatusAvailable   CRVInfoStatus = "available"
	CRVInfoStatusUnavailable CRVInfoStatus = "unavailable"
)

type CRVInfoRequest struct {
	Data CRVInfoRequestData `json:"data"`
}

type CRVInfoRequestData struct {
	Release  ReleaseRef  `json:"release"`
	Stemcell StemcellRef `json:"stemcell"`
}

type CRVInfoResponse struct {
	Data CRVInfoResponseData `json:"data"`
}

type CRVInfoResponseData struct {
	Status   CRVInfoStatus               `json:"status"`
	Release  ReleaseRef                  `json:"release,omitempty"`
	Stemcell StemcellRef                 `json:"stemcell,omitempty"`
	Tarball  CRVInfoResponseDataCompiled `json:"tarball,omitempty"`
}

type CRVInfoResponseDataCompiled struct {
	URL       string    `json:"url"`
	Checksums Checksums `json:"checksums"`
}
