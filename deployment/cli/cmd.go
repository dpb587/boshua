package cli

import (
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
)

type CmdOpts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`
}

type Cmd struct {
	PatchManifestCmd PatchManifestCmd `command:"patch-manifest" description:"For patching a manifest to refer to compiled releases"`
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{}

	cmdOpts := &CmdOpts{
		AppOpts: app,
	}

	cmd.PatchManifestCmd.CmdOpts = cmdOpts

	return cmd
}
