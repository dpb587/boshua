package models

type ScheduleRequest struct {
	Data ScheduleRequestData `json:"data"`
}

type ScheduleRequestData struct {
	ReleaseVersionRef  ReleaseVersionRef  `json:"release_version_ref"`
	StemcellVersionRef StemcellVersionRef `json:"stemcell_version_ref"`
}

type ScheduleResponse struct {
	Status string `json:"status"`
}
