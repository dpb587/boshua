package boshuaV2

import (
	"github.com/dpb587/boshua/stemcellversion"
)

type filterResponse struct {
	Stemcells []stemcellversion.Artifact `json:"stemcells"`
}
