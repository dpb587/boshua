package cli

import (
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
)

type Cmd struct {
	PropertiesCmd PropertiesCmd `command:"properties" description:"Show the job properties"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.PropertiesCmd.Execute(extra)
}

type CmdOpts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{}

	cmdOpts := &CmdOpts{
		AppOpts: app,
	}

	cmd.PropertiesCmd.CmdOpts = cmdOpts

	return cmd
}
