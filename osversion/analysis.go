package osversion

import "github.com/dpb587/boshua/analysis"

var _ analysis.Subject = &Artifact{}

func (s Artifact) SupportedAnalyzers() []string {
	return []string{
		"osimagechecksums.v1",
		"osimagefilestat.v1",
		"osmanifest.v1",
	}
}
