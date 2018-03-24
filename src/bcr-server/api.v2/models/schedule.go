package models

type ScheduleRequest struct {
	Data ScheduleRequestData `json:"data"`
}

type ScheduleRequestData struct {
	Release  ReleaseRef  `json:"release"`
	Stemcell StemcellRef `json:"stemcell"`
}

type ScheduleResponse struct {
	Status string `json:"status"`
}
