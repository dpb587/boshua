package datastore

import (
	"github.com/dpb587/boshua/cli/client/cmd/compiledrelease/datastore/opts"
	compiledreleaseopts "github.com/dpb587/boshua/cli/client/cmd/compiledrelease/opts"
	cmdopts "github.com/dpb587/boshua/cli/client/cmd/opts"
	compiledreleaseversiondatastore "github.com/dpb587/boshua/compiledreleaseversion/datastore"
)

type Cmd struct {
	*opts.Opts

	FilterCmd FilterCmd `command:"filter" description:"For filtering results"`
}

type CmdOpts struct {
	AppOpts             *cmdopts.Opts `no-flag:"true"`
	CompiledReleaseOpts *compiledreleaseopts.Opts
	DatastoreOpts       *opts.Opts
}

func (o *CmdOpts) getDatastore() (compiledreleaseversiondatastore.Index, error) {
	return o.AppOpts.GetCompiledReleaseIndex(o.DatastoreOpts.Datastore)
}

func New(app *cmdopts.Opts, compiledrelease *compiledreleaseopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmdOpts := &CmdOpts{
		AppOpts:             app,
		CompiledReleaseOpts: compiledrelease,
		DatastoreOpts:       cmd.Opts,
	}

	cmd.FilterCmd.CmdOpts = cmdOpts

	return cmd
}
