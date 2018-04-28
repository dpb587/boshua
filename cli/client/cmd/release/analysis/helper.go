package analysis

import (
	"fmt"
	"os"

	"github.com/dpb587/boshua/api/v2/models/analysis"
	"github.com/dpb587/boshua/api/v2/models/scheduler"
	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/releaseversion"
)

func (o *CmdOpts) getAnalysis() (*analysis.GETAnalysisResponse, error) {
	client := o.AppOpts.GetClient()

	ref := releaseversion.Reference{
		Name:      o.ReleaseOpts.Release.Name,
		Version:   o.ReleaseOpts.Release.Version,
		Checksums: checksum.ImmutableChecksums{o.ReleaseOpts.ReleaseChecksum.ImmutableChecksum},
	}
	analyzer := o.AnalysisOpts.Analyzer

	if o.AnalysisOpts.NoWait {
		return client.GetReleaseVersionAnalysis(ref, analyzer)
	}

	return client.RequireReleaseVersionAnalysis(
		ref,
		analyzer,
		func(task scheduler.TaskStatus) {
			if o.AppOpts.Quiet {
				return
			}

			fmt.Fprintf(os.Stderr, "analysis status: %s\n", task.Status)
		},
	)
}
