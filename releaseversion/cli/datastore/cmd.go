package datastore

import (
	cmdopts "github.com/dpb587/boshua/main/boshua/cmd/opts"
	"github.com/dpb587/boshua/releaseversion/cli/datastore/opts"
	releaseopts "github.com/dpb587/boshua/releaseversion/cli/opts"
)

type Cmd struct {
	*opts.Opts

	FilterCmd FilterCmd `command:"filter" description:"For filtering results"`
	LabelsCmd LabelsCmd `command:"labels" description:"For listing all labels"`
}

type CmdOpts struct {
	AppOpts       *cmdopts.Opts     `no-flag:"true"`
	ReleaseOpts   *releaseopts.Opts `no-flag:"true"`
	DatastoreOpts *opts.Opts        `no-flag:"true"`
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
	cmd.LabelsCmd.CmdOpts = cmdOpts

	return cmd
}
