package datastore

import (
	cmdopts "github.com/dpb587/boshua/main/boshua/cmd/opts"
	"github.com/dpb587/boshua/releaseversion/compilation/cli/datastore/opts"
	compiledreleaseopts "github.com/dpb587/boshua/releaseversion/compilation/cli/opts"
)

type Cmd struct {
	*opts.Opts

	FilterCmd FilterCmd `command:"filter" description:"For filtering results"`
	StoreCmd  StoreCmd  `command:"store" description:"For storing an artifact"`
}

type CmdOpts struct {
	AppOpts             *cmdopts.Opts `no-flag:"true"`
	CompiledReleaseOpts *compiledreleaseopts.Opts
}

func New(app *cmdopts.Opts, compiledrelease *compiledreleaseopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmdOpts := &CmdOpts{
		AppOpts:             app,
		CompiledReleaseOpts: compiledrelease,
	}

	cmd.FilterCmd.CmdOpts = cmdOpts
	cmd.StoreCmd.CmdOpts = cmdOpts

	return cmd
}
