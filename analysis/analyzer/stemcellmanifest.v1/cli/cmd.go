package cli

import (
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
)

type Cmd struct {
	ContentsCmd ContentsCmd `command:"contents" description:"Show the original contents of stemcell.MF"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.ContentsCmd.Execute(extra)
}

type CmdOpts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{}

	cmdOpts := &CmdOpts{
		AppOpts: app,
	}

	cmd.ContentsCmd.CmdOpts = cmdOpts

	return cmd
}
