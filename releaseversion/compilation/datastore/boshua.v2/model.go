package boshuaV2

import (
	"github.com/dpb587/boshua/releaseversion/compilation"
)

type filterResponse struct {
	Release filterReleaseResponse `json:"release"`
}

type filterReleaseResponse struct {
	Compilation compilation.Artifact `json:"compilation"`
}
