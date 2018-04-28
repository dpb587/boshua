package analysis

import (
	"fmt"
	"os"

	"github.com/dpb587/boshua/api/v2/models/analysis"
	"github.com/dpb587/boshua/api/v2/models/scheduler"
	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/cli/client/cmd/analysisutil/opts"
	compiledreleaseopts "github.com/dpb587/boshua/cli/client/cmd/compiledrelease/opts"
	cmdopts "github.com/dpb587/boshua/cli/client/cmd/opts"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
)

type Cmd struct {
	*opts.Opts

	MetalinkCmd MetalinkCmd `command:"metalink" description:"For showing a metalink of the analysis"`
	ResultsCmd  ResultsCmd  `command:"results" description:"For showing the results of an analysis"`
}

type CmdOpts struct {
	AppOpts             *cmdopts.Opts `no-flag:"true"`
	CompiledReleaseOpts *compiledreleaseopts.Opts
	AnalysisOpts        *opts.Opts
}

func (o *CmdOpts) getAnalysis() (*analysis.GETAnalysisResponse, error) {
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

			fmt.Fprintf(os.Stderr, "analysis status: %s\n", task.Status)
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

	cmd.MetalinkCmd.CmdOpts = cmdOpts
	cmd.ResultsCmd.CmdOpts = cmdOpts

	return cmd
}
