package datastore

import (
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	"github.com/dpb587/boshua/releaseversion/cli/datastore/opts"
	releaseopts "github.com/dpb587/boshua/releaseversion/cli/opts"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
)

type Cmd struct {
	*opts.Opts

	FilterCmd FilterCmd `command:"filter" description:"For filtering results"`
}

type CmdOpts struct {
	AppOpts       *cmdopts.Opts     `no-flag:"true"`
	ReleaseOpts   *releaseopts.Opts `no-flag:"true"`
	DatastoreOpts *opts.Opts
}

func (o *CmdOpts) getDatastore() (releaseversiondatastore.Index, error) {
	return o.AppOpts.GetReleaseIndex(o.DatastoreOpts.Datastore)
}

func New(app *cmdopts.Opts, release *releaseopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmdOpts := &CmdOpts{
		AppOpts:       app,
		ReleaseOpts:   release,
		DatastoreOpts: cmd.Opts,
	}

	cmd.FilterCmd.CmdOpts = cmdOpts

	return cmd
}
