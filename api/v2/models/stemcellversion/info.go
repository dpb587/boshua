package stemcellversion

import (
	"github.com/dpb587/metalink"
)

type InfoResponse struct {
	Data InfoResponseData `json:"data"`
}

type InfoResponseData struct {
	Reference Reference     `json:"reference"`
	Artifact  metalink.File `json:"file"`
}
