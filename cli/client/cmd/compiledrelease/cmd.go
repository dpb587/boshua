package compiledrelease

import (
	"github.com/dpb587/bosh-compiled-releases/cli/client/cmd/compiledrelease/opts"
	cmdopts "github.com/dpb587/bosh-compiled-releases/cli/client/cmd/opts"
)

type CmdOpts struct {
	AppOpts             *cmdopts.Opts `no-flag:"true"`
	CompiledReleaseOpts *opts.Opts
}

type Cmd struct {
	*opts.Opts

	DownloadCmd DownloadCmd `command:"download" description:"For downloading a compiled release tarball"`
	MetalinkCmd MetalinkCmd `command:"metalink" description:"For showing a metalink of the compiled release"`
	OpsFileCmd  OpsFileCmd  `command:"ops-file" description:"For showing a deployment manifest ops file for the compiled release"`
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmdOpts := &CmdOpts{
		AppOpts:             app,
		CompiledReleaseOpts: cmd.Opts,
	}

	cmd.MetalinkCmd.CmdOpts = cmdOpts
	cmd.DownloadCmd.CmdOpts = cmdOpts
	cmd.OpsFileCmd.CmdOpts = cmdOpts

	return cmd
}