package datastore

import (
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	"github.com/dpb587/boshua/stemcellversion/cli/datastore/opts"
	stemcellopts "github.com/dpb587/boshua/stemcellversion/cli/opts"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
)

type Cmd struct {
	*opts.Opts

	FilterCmd        FilterCmd        `command:"filter" description:"For filtering results"`
	StoreAnalysisCmd StoreAnalysisCmd `command:"store-analysis" description:"For storing analysis results"`
}

type CmdOpts struct {
	AppOpts       *cmdopts.Opts      `no-flag:"true"`
	StemcellOpts  *stemcellopts.Opts `no-flag:"true"`
	DatastoreOpts *opts.Opts
}

func (o *CmdOpts) getDatastore() (stemcellversiondatastore.Index, error) {
	return o.AppOpts.GetStemcellIndex(o.DatastoreOpts.Datastore)
}

func New(app *cmdopts.Opts, stemcell *stemcellopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmdOpts := &CmdOpts{
		AppOpts:       app,
		StemcellOpts:  stemcell,
		DatastoreOpts: cmd.Opts,
	}

	cmd.StoreAnalysisCmd.CmdOpts = cmdOpts
	cmd.FilterCmd.CmdOpts = cmdOpts

	return cmd
}
