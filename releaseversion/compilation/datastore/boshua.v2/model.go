package boshuaV2

import (
	"github.com/dpb587/metalink"
)

type filterResponse struct {
	Releases []filterReleaseResponse `json:"releases"`
}

type filterReleaseResponse struct {
	Name         string                             `json:"name"`
	Version      string                             `json:"version"`
	Compilations []filterReleaseCompilationResponse `json:"compilations"`
}

type filterReleaseCompilationResponse struct {
	OS      string        `json:"os"`
	Version string        `json:"version"`
	Labels  []string      `json:"labels"`
	Tarball metalink.File `json:"tarball"`
}
