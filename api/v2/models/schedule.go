package models

type ScheduleRequest struct {
	Data ScheduleRequestData `json:"data"`
}

type ScheduleRequestData struct {
	ReleaseVersionRef ReleaseVersionRef `json:"release_version_ref"`
	OSVersionRef      OSVersionRef      `json:"os_version_ref"`
}

type ScheduleResponse struct {
	Status string `json:"status"`
}
