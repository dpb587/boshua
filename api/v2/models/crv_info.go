package models

import (
	"github.com/dpb587/metalink"
)

type CRVInfoRequest struct {
	Data CRVInfoRequestData `json:"data"`
}

type CRVInfoRequestData struct {
	ReleaseVersionRef  ReleaseVersionRef  `json:"release_version_ref"`
	OSVersionRef OSVersionRef `json:"os_version_ref"`
}

type CRVInfoResponse struct {
	Data CRVInfoResponseData `json:"data"`
}

type CRVInfoResponseData struct {
	ReleaseVersionRef  ReleaseVersionRef  `json:"release_version_ref,omitempty"`
	OSVersionRef OSVersionRef `json:"os_version_ref,omitempty"`
	Artifact           metalink.File      `json:"artifact,omitempty"`
}
