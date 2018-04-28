package analysis

import (
	cmdopts "github.com/dpb587/boshua/cli/client/cmd/opts"
	"github.com/dpb587/boshua/cli/client/cmd/release/analysis/opts"
	releaseopts "github.com/dpb587/boshua/cli/client/cmd/release/opts"
)

type CmdOpts struct {
	AppOpts      *cmdopts.Opts `no-flag:"true"`
	ReleaseOpts  *releaseopts.Opts
	AnalysisOpts *opts.Opts
}

type Cmd struct {
	*opts.Opts

	MetalinkCmd MetalinkCmd `command:"metalink" description:"For showing a metalink of the compiled release"`
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

	return cmd
}
