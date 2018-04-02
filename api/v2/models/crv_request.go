package models

type CRVRequestRequest struct {
	Data CRVRequestRequestData `json:"data"`
}

type CRVRequestRequestData struct {
	Release  ReleaseRef  `json:"release"`
	Stemcell StemcellRef `json:"stemcell"`
}

type CRVRequestResponse struct {
	Complete bool   `json:"complete"`
	Status   string `json:"status"`
}
