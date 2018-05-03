package analysis

import (
	"github.com/dpb587/metalink"
)

type GETInfoResponse struct {
	Data GETInfoResponseData `json:"data"`
}

type GETInfoResponseData struct {
	Artifact metalink.File `json:"artifact"`
}
