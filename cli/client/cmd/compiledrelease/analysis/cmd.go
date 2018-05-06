package analysis

import (
	"fmt"
	"os"
	"time"

	"github.com/dpb587/boshua/api/v2/models/analysis"
	"github.com/dpb587/boshua/api/v2/models/scheduler"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/boshua/cli/client/cmd/analysisutil/opts"
	compiledreleaseopts "github.com/dpb587/boshua/cli/client/cmd/compiledrelease/opts"
	cmdopts "github.com/dpb587/boshua/cli/client/cmd/opts"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
)

type Cmd struct {
	*opts.Opts

	ArtifactCmd ArtifactCmd `command:"artifact" description:"For showing the analysis artifact"`
	ResultsCmd  ResultsCmd  `command:"results" description:"For showing the results of an analysis"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.ArtifactCmd.Execute(extra)
}

type CmdOpts struct {
	AppOpts             *cmdopts.Opts `no-flag:"true"`
	CompiledReleaseOpts *compiledreleaseopts.Opts
	AnalysisOpts        *opts.Opts
}

func (o *CmdOpts) getAnalysis() (*analysis.GETInfoResponse, error) {
	client := o.AppOpts.GetClient()

	releaseVersionRef := releaseversion.Reference{
		Name:      o.CompiledReleaseOpts.Release.Name,
		Version:   o.CompiledReleaseOpts.Release.Version,
		Checksums: checksum.ImmutableChecksums{o.CompiledReleaseOpts.ReleaseChecksum.ImmutableChecksum},
	}
	osVersionRef := osversion.Reference{
		Name:    o.CompiledReleaseOpts.OS.Name,
		Version: o.CompiledReleaseOpts.OS.Version,
	}
	analyzer := o.AnalysisOpts.Analyzer

	if o.AnalysisOpts.NoWait {
		return client.GetCompiledReleaseVersionAnalysis(releaseVersionRef, osVersionRef, analyzer)
	}

	return client.RequireCompiledReleaseVersionAnalysis(
		releaseVersionRef,
		osVersionRef,
		analyzer,
		func(task scheduler.TaskStatus) {
			if o.AppOpts.Quiet {
				return
			}

			fmt.Fprintf(os.Stderr, "boshua | %s | requesting compiled release analysis: %s/%s: %s/%s: %s: task is %s\n", time.Now().Format("15:04:05"), osVersionRef.Name, osVersionRef.Version, releaseVersionRef.Name, releaseVersionRef.Version, analyzer, task.Status)
		},
	)
}

func New(app *cmdopts.Opts, compiledrelease *compiledreleaseopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmdOpts := &CmdOpts{
		AppOpts:             app,
		CompiledReleaseOpts: compiledrelease,
		AnalysisOpts:        cmd.Opts,
	}

	cmd.ArtifactCmd.CmdOpts = cmdOpts
	cmd.ResultsCmd.CmdOpts = cmdOpts

	return cmd
}
