package cli

import (
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
)

type CmdOpts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`
}

type Cmd struct {
	UseCompiledReleasesCmd UseCompiledReleasesCmd `command:"use-compiled-releases" description:"For patching a manifest to refer to compiled releases"`
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{}

	cmdOpts := &CmdOpts{
		AppOpts: app,
	}

	cmd.UseCompiledReleasesCmd.CmdOpts = cmdOpts

	return cmd
}
