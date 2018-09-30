package cli

import (
	"github.com/dpb587/boshua/config/provider/setter"
	cmdopts "github.com/dpb587/boshua/main/boshua/cmd/opts"
	"github.com/dpb587/boshua/pivnetfile/cli/analysis"
	"github.com/dpb587/boshua/pivnetfile/cli/datastore"
	"github.com/dpb587/boshua/pivnetfile/cli/opts"
)

type CmdOpts struct {
	PivnetFileOpts *opts.Opts
}

type Cmd struct {
	setter.AppConfig `no-flag:"true"`
	*opts.Opts

	AnalysisCmd    *analysis.Cmd    `command:"analysis" description:"For analyzing the pivnet file" subcommands-optional:"true"`
	DatastoreCmd   *datastore.Cmd   `command:"datastore" description:"For interacting with pivnet file datastores"`

	AnalyzersCmd     AnalyzersCmd     `command:"analyzers" description:"For showing the supported analyzers"`
	ArtifactCmd      ArtifactCmd      `command:"artifact" description:"For showing the pivnet file artifact"`
	DownloadCmd      DownloadCmd      `command:"download" description:"For downloading the pivnet file locally"`
}

func (c *Cmd) Execute(extra []string) error {
	c.ArtifactCmd.SetConfig(c.AppConfig.Config)
	return c.ArtifactCmd.Execute(extra)
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmd.AnalysisCmd = analysis.New(app, cmd.Opts)
	cmd.DatastoreCmd = datastore.New(app, cmd.Opts)

	cmdOpts := &CmdOpts{
		PivnetFileOpts: cmd.Opts,
	}

	cmd.AnalyzersCmd.CmdOpts = cmdOpts
	cmd.ArtifactCmd.CmdOpts = cmdOpts
	cmd.DownloadCmd.CmdOpts = cmdOpts

	return cmd
}
