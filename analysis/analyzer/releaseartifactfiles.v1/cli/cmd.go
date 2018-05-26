package cli

import (
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
)

type Cmd struct {
	LsCmd        LsCmd        `command:"ls" description:"Show an ls-style list of files"`
	Sha1sumCmd   Sha1sumCmd   `command:"sha1sum" alias:"shasum" description:"Show sha1 checksums"`
	Sha256sumCmd Sha256sumCmd `command:"sha256sum" description:"Show sha256 checksums"`
	Sha512sumCmd Sha512sumCmd `command:"sha512sum" description:"Show sha512 checksums"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.LsCmd.Execute(extra)
}

type CmdOpts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{}

	cmdOpts := &CmdOpts{
		AppOpts: app,
	}

	cmd.LsCmd.CmdOpts = cmdOpts
	cmd.Sha1sumCmd.CmdOpts = cmdOpts
	cmd.Sha256sumCmd.CmdOpts = cmdOpts
	cmd.Sha512sumCmd.CmdOpts = cmdOpts

	return cmd
}
