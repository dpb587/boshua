package stemcell

import (
	cmdopts "github.com/dpb587/boshua/cli/client/cmd/opts"
	"github.com/dpb587/boshua/cli/client/cmd/stemcell/analysis"
	"github.com/dpb587/boshua/cli/client/cmd/stemcell/opts"
)

type CmdOpts struct {
	AppOpts      *cmdopts.Opts `no-flag:"true"`
	StemcellOpts *opts.Opts
}

type Cmd struct {
	*opts.Opts

	AnalysisCmd *analysis.Cmd `command:"analysis" description:"For analyzing the stemcell artifact"`

	ArtifactCmd ArtifactCmd `command:"artifact" description:"For showing the stemcell artifact"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.ArtifactCmd.Execute(extra)
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmd.AnalysisCmd = analysis.New(app, cmd.Opts)

	cmdOpts := &CmdOpts{
		AppOpts:      app,
		StemcellOpts: cmd.Opts,
	}

	cmd.ArtifactCmd.CmdOpts = cmdOpts

	return cmd
}
