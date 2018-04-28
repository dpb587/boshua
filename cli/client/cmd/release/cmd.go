package release

import (
	cmdopts "github.com/dpb587/boshua/cli/client/cmd/opts"
	"github.com/dpb587/boshua/cli/client/cmd/release/analysis"
	"github.com/dpb587/boshua/cli/client/cmd/release/opts"
)

type CmdOpts struct {
	AppOpts     *cmdopts.Opts `no-flag:"true"`
	ReleaseOpts *opts.Opts
}

type Cmd struct {
	*opts.Opts

	AnalysisCmd *analysis.Cmd `command:"analysis" description:"For analyzing artifacts"`
	MetalinkCmd MetalinkCmd   `command:"metalink" description:"For showing a metalink of the compiled release"`
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmd.AnalysisCmd = analysis.New(app, cmd.Opts)

	cmdOpts := &CmdOpts{
		AppOpts:     app,
		ReleaseOpts: cmd.Opts,
	}

	cmd.MetalinkCmd.CmdOpts = cmdOpts

	return cmd
}
