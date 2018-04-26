package stemcellversion

import "github.com/dpb587/boshua/analysis"

var _ analysis.Subject = &Artifact{}

func (s Artifact) SupportedAnalyzers() []string {
	return []string{
		"stemcellimagechecksums.v1",
		"stemcellimagefilestat.v1",
		"stemcellmanifest.v1",
	}
}
