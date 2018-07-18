package analysis

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon/opts"
	"github.com/dpb587/boshua/analysis/cli/cliutil"
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	compiledreleaseopts "github.com/dpb587/boshua/releaseversion/compilation/cli/opts"
)

type Cmd struct {
	*opts.Opts

	ArtifactCmd     ArtifactCmd     `command:"artifact" description:"For showing the analysis artifact"`
	ResultsCmd      ResultsCmd      `command:"results" description:"For showing the results of an analysis"`
	StoreResultsCmd StoreResultsCmd `command:"store-results" description:"For storing the results of an analysis"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.ArtifactCmd.Execute(extra)
}

type CmdOpts struct {
	AppOpts             *cmdopts.Opts `no-flag:"true"`
	CompiledReleaseOpts *compiledreleaseopts.Opts
	AnalysisOpts        *opts.Opts
}

func (o *CmdOpts) getAnalysis() (analysis.Artifact, error) {
	return cliutil.LoadAnalysis(
		o.AppOpts.GetAnalysisIndex,
		func() (analysis.Subject, error) {
			return o.CompiledReleaseOpts.Artifact()
		},
		o.AnalysisOpts,
		o.AppOpts.GetScheduler,
		append(
			append(
				append(
					[]string{"release"},
					o.CompiledReleaseOpts.ReleaseOpts.Opts()...,
				),
				"compilation",
			),
			fmt.Sprintf("--os=%s", o.CompiledReleaseOpts.OS.String()),
		),
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
	cmd.StoreResultsCmd.CmdOpts = cmdOpts

	return cmd
}
