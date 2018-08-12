package cli

import (
	cmdopts "github.com/dpb587/boshua/main/boshua/cmd/opts"
	"github.com/dpb587/boshua/releaseversion/cli/analysis"
	"github.com/dpb587/boshua/releaseversion/cli/datastore"
	"github.com/dpb587/boshua/releaseversion/cli/opts"
	compilation "github.com/dpb587/boshua/releaseversion/compilation/cli"
)

type CmdOpts struct {
	AppOpts     *cmdopts.Opts `no-flag:"true"`
	ReleaseOpts *opts.Opts
}

type Cmd struct {
	*opts.Opts

	AnalysisCmd    *analysis.Cmd    `command:"analysis" description:"For analyzing the release artifact" subcommands-optional:"true"`
	DatastoreCmd   *datastore.Cmd   `command:"datastore" description:"For interacting with release datastores"`
	CompilationCmd *compilation.Cmd `command:"compilation" description:"For working with compiled releases" subcommands-optional:"true"`

	AnalyzersCmd     AnalyzersCmd     `command:"analyzers" description:"For showing the supported analyzers"`
	ArtifactCmd      ArtifactCmd      `command:"artifact" description:"For showing the release artifact"`
	UploadReleaseCmd UploadReleaseCmd `command:"upload-release" description:"For uploading the release to BOSH"`
	DownloadCmd      DownloadCmd      `command:"download" description:"For downloading the release locally"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.ArtifactCmd.Execute(extra)
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{
			AppOpts: app,
		},
	}

	cmd.AnalysisCmd = analysis.New(app, cmd.Opts)
	cmd.DatastoreCmd = datastore.New(app, cmd.Opts)
	cmd.CompilationCmd = compilation.New(app, cmd.Opts)

	cmdOpts := &CmdOpts{
		AppOpts:     app,
		ReleaseOpts: cmd.Opts,
	}

	cmd.AnalyzersCmd.CmdOpts = cmdOpts
	cmd.ArtifactCmd.CmdOpts = cmdOpts
	cmd.UploadReleaseCmd.CmdOpts = cmdOpts
	cmd.DownloadCmd.CmdOpts = cmdOpts

	return cmd
}
