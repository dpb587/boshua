package analysis

import (
	"github.com/dpb587/boshua/cli/client/cmd/analysis/opts"
	cmdopts "github.com/dpb587/boshua/cli/client/cmd/opts"
)

type CmdOpts struct {
	AppOpts      *cmdopts.Opts `no-flag:"true"`
	AnalysisOpts *opts.Opts
}

type Cmd struct {
	*opts.Opts

	GenerateCmd GenerateCmd `command:"generate" description:"For generating an analysis"`
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmdOpts := &CmdOpts{
		AppOpts:      app,
		AnalysisOpts: cmd.Opts,
	}

	cmd.GenerateCmd.CmdOpts = cmdOpts

	return cmd
}
