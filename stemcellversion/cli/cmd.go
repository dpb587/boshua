package cli

import (
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	"github.com/dpb587/boshua/stemcellversion/cli/analysis"
	"github.com/dpb587/boshua/stemcellversion/cli/datastore"
	"github.com/dpb587/boshua/stemcellversion/cli/opts"
)

type CmdOpts struct {
	AppOpts      *cmdopts.Opts `no-flag:"true"`
	StemcellOpts *opts.Opts
}

type Cmd struct {
	*opts.Opts

	AnalysisCmd  *analysis.Cmd  `command:"analysis" description:"For analyzing the stemcell artifact" subcommands-optional:"true"`
	DatastoreCmd *datastore.Cmd `command:"datastore" description:"For interacting with release datastores"`

	AnalyzersCmd      AnalyzersCmd      `command:"analyzers" description:"For showing the supported analyzers"`
	ArtifactCmd       ArtifactCmd       `command:"artifact" description:"For showing the stemcell artifact"`
	UploadStemcellCmd UploadStemcellCmd `command:"upload-stemcell" description:"For uploading the stemcell to BOSH"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.ArtifactCmd.Execute(extra)
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmd.AnalysisCmd = analysis.New(app, cmd.Opts)
	cmd.DatastoreCmd = datastore.New(app, cmd.Opts)

	cmdOpts := &CmdOpts{
		AppOpts:      app,
		StemcellOpts: cmd.Opts,
	}

	cmd.AnalyzersCmd.CmdOpts = cmdOpts
	cmd.ArtifactCmd.CmdOpts = cmdOpts
	cmd.UploadStemcellCmd.CmdOpts = cmdOpts

	return cmd
}
