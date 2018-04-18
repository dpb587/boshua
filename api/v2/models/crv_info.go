package models

import "time"

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
	Release  ReleaseRef                  `json:"release,omitempty"`
	Stemcell StemcellRef                 `json:"stemcell,omitempty"`
	Tarball  CRVInfoResponseDataCompiled `json:"tarball,omitempty"`
}

type CRVInfoResponseDataCompiled struct {
	URLs      []string   `json:"urls"`
	Size      *uint64    `json:"size,omitempty"`
	Published *time.Time `json:"published,omitempty"`
	Checksums Checksums  `json:"checksums"`
}
