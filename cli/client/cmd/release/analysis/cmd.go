package analysis

import (
	"fmt"
	"os"

	"github.com/dpb587/boshua/api/v2/models/analysis"
	"github.com/dpb587/boshua/api/v2/models/scheduler"
	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/cli/client/cmd/analysisutil/opts"
	cmdopts "github.com/dpb587/boshua/cli/client/cmd/opts"
	releaseopts "github.com/dpb587/boshua/cli/client/cmd/release/opts"
	"github.com/dpb587/boshua/releaseversion"
)

type Cmd struct {
	*opts.Opts

	MetalinkCmd MetalinkCmd `command:"metalink" description:"For showing a metalink of the analysis"`
	ResultsCmd  ResultsCmd  `command:"results" description:"For showing the results of an analysis"`
}

type CmdOpts struct {
	AppOpts      *cmdopts.Opts `no-flag:"true"`
	ReleaseOpts  *releaseopts.Opts
	AnalysisOpts *opts.Opts
}

func (o *CmdOpts) getAnalysis() (*analysis.GETInfoResponse, error) {
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

func New(app *cmdopts.Opts, release *releaseopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmdOpts := &CmdOpts{
		AppOpts:      app,
		ReleaseOpts:  release,
		AnalysisOpts: cmd.Opts,
	}

	cmd.MetalinkCmd.CmdOpts = cmdOpts
	cmd.ResultsCmd.CmdOpts = cmdOpts

	return cmd
}
