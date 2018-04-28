package analysis

import (
	"github.com/dpb587/boshua/api/v2/models/analysis"
	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/releaseversion"
)

func (o *CmdOpts) getAnalysis() (*analysis.GETAnalysisResponse, error) {
	client := o.AppOpts.GetClient()

	var get func(releaseversion.Reference, string) (*analysis.GETAnalysisResponse, error) = client.GetReleaseVersionAnalysis

	if o.AnalysisOpts.RequestAndWait {
		get = client.RequireReleaseVersionAnalysis
	}

	return get(
		releaseversion.Reference{
			Name:      o.ReleaseOpts.Release.Name,
			Version:   o.ReleaseOpts.Release.Version,
			Checksums: checksum.ImmutableChecksums{o.ReleaseOpts.ReleaseChecksum.ImmutableChecksum},
		},
		o.AnalysisOpts.Analyzer,
	)
}
