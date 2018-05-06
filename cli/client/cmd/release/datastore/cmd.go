package datastore

import (
	"time"

	cmdopts "github.com/dpb587/boshua/cli/client/cmd/opts"
	"github.com/dpb587/boshua/cli/client/cmd/release/datastore/opts"
	releaseopts "github.com/dpb587/boshua/cli/client/cmd/release/opts"
	"github.com/dpb587/boshua/datastore/git"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/boshreleasedir"
	"github.com/sirupsen/logrus"
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
	// client := o.AppOpts.GetClient()

	return boshreleasedir.New(boshreleasedir.Config{
		RepositoryConfig: git.RepositoryConfig{
			Repository:   "git+https://github.com/dpb587/ssoca-bosh-release.git",
			LocalPath:    "/Users/dpb587/Projects/src/github.com/dpb587/ssoca-bosh-release",
			PullInterval: 30 * time.Minute,
		},
	}, logrus.New()), nil
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
