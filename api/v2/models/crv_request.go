package models

type CRVRequestRequest struct {
	Data CRVRequestRequestData `json:"data"`
}

type CRVRequestRequestData struct {
	Release  ReleaseRef  `json:"release"`
	Stemcell StemcellRef `json:"stemcell"`
}

type CRVRequestResponse struct {
	Status string `json:"status"`
}
