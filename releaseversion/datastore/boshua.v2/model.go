package boshuaV2

import (
	"github.com/dpb587/boshua/releaseversion"
)

type filterResponse struct {
	Releases []releaseversion.Artifact `json:"releases"`
}
