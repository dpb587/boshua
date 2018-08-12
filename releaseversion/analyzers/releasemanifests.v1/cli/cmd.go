package cli

import (
	cmdopts "github.com/dpb587/boshua/main/boshua/cmd/opts"
)

type Cmd struct {
	PropertiesCmd PropertiesCmd `command:"properties" description:"Show the job properties"`
	SpecCmd       SpecCmd       `command:"spec" description:"Show the job or release manifests"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.SpecCmd.Execute(extra)
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
	cmd.SpecCmd.CmdOpts = cmdOpts

	return cmd
}
