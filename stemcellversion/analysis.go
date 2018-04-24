package stemcellversion

import "github.com/dpb587/boshua/analysis"

var _ analysis.Subject = &Subject{}

func (s Subject) SupportedAnalyzers() []string {
	return []string{
		"stemcellimagechecksums.v1",
		"stemcellimagefilestat.v1",
		"stemcellmanifest.v1",
	}
}
