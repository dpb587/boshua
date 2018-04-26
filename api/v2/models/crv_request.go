package models

type CRVRequestRequest struct {
	Data CRVRequestRequestData `json:"data"`
}

type CRVRequestRequestData struct {
	ReleaseVersionRef  ReleaseVersionRef  `json:"release_version_ref"`
	OSVersionRef OSVersionRef `json:"stemcell_version_ref"`
}

type CRVRequestResponse struct {
	Complete bool   `json:"complete"`
	Status   string `json:"status"`
}
