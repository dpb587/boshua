package analyzer

import (
	cmdopts "github.com/dpb587/boshua/cli/client/cmd/opts"
)

type Cmd struct {
	PackagesCmd PackagesCmd `command:"packages" description:"Show a simple list of package versions"`
	ContentsCmd ContentsCmd `command:"contents" description:"Show the original contents of packages.txt"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.PackagesCmd.Execute(extra)
}

type CmdOpts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{}

	cmdOpts := &CmdOpts{
		AppOpts: app,
	}

	cmd.PackagesCmd.CmdOpts = cmdOpts
	cmd.ContentsCmd.CmdOpts = cmdOpts

	return cmd
}
