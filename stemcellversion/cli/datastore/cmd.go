package datastore

import (
	cmdopts "github.com/dpb587/boshua/main/boshua/cmd/opts"
	stemcellopts "github.com/dpb587/boshua/stemcellversion/cli/opts"
)

type Cmd struct {
	FilterCmd FilterCmd `command:"filter" description:"For filtering results"`
}

type CmdOpts struct {
	AppOpts      *cmdopts.Opts      `no-flag:"true"`
	StemcellOpts *stemcellopts.Opts `no-flag:"true"`
}

func New(app *cmdopts.Opts, stemcell *stemcellopts.Opts) *Cmd {
	cmd := &Cmd{}

	cmdOpts := &CmdOpts{
		AppOpts:      app,
		StemcellOpts: stemcell,
	}

	cmd.FilterCmd.CmdOpts = cmdOpts

	return cmd
}
