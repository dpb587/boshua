package analysis

import (
	"fmt"
	"os"
	"time"

	"github.com/dpb587/boshua/api/v2/models/analysis"
	"github.com/dpb587/boshua/api/v2/models/scheduler"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/boshua/cli/client/cmd/analysisutil/opts"
	cmdopts "github.com/dpb587/boshua/cli/client/cmd/opts"
	releaseopts "github.com/dpb587/boshua/cli/client/cmd/release/opts"
	"github.com/dpb587/boshua/releaseversion"
)

type Cmd struct {
	*opts.Opts

	ArtifactCmd ArtifactCmd `command:"metalink" description:"For showing the analysis artifact"`
	ResultsCmd  ResultsCmd  `command:"results" description:"For showing the results of an analysis"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.ArtifactCmd.Execute(extra)
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

			fmt.Fprintf(os.Stderr, "boshua | %s | requesting release analysis: %s/%s: %s: task is %s\n", time.Now().Format("15:04:05"), ref.Name, ref.Version, analyzer, task.Status)
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

	cmd.ArtifactCmd.CmdOpts = cmdOpts
	cmd.ResultsCmd.CmdOpts = cmdOpts

	return cmd
}
