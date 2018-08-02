package cli

import (
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	releaseopts "github.com/dpb587/boshua/releaseversion/cli/opts"
	"github.com/dpb587/boshua/releaseversion/compilation/cli/analysis"
	"github.com/dpb587/boshua/releaseversion/compilation/cli/datastore"
	"github.com/dpb587/boshua/releaseversion/compilation/cli/opts"
)

type CmdOpts struct {
	AppOpts             *cmdopts.Opts `no-flag:"true"`
	CompiledReleaseOpts *opts.Opts
}

type Cmd struct {
	*opts.Opts

	AnalysisCmd  *analysis.Cmd  `command:"analysis" description:"For analyzing artifacts" subcommands-optional:"true" subcommands-optional:"true"`
	DatastoreCmd *datastore.Cmd `command:"datastore" description:"For interacting with compiled release datastores"`

	AnalyzersCmd     AnalyzersCmd     `command:"analyzers" description:"For showing the supported analyzers"`
	ArtifactCmd      ArtifactCmd      `command:"artifact" description:"For showing the compiled release artifact"`
	OpsFileCmd       OpsFileCmd       `command:"ops-file" description:"For showing a deployment manifest ops file for the compiled release"`
	UploadReleaseCmd UploadReleaseCmd `command:"upload-release" description:"For uploading the compiled release to BOSH"`
	DownloadCmd      DownloadCmd      `command:"download" description:"For downloading the compiled release locally"`
	ExportReleaseCmd ExportReleaseCmd `command:"export-release" description:"For exporting a compiled release from BOSH"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.ArtifactCmd.Execute(extra)
}

func New(app *cmdopts.Opts, release *releaseopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{
			AppOpts:     app,
			ReleaseOpts: release,
		},
	}

	cmd.AnalysisCmd = analysis.New(app, cmd.Opts)
	cmd.DatastoreCmd = datastore.New(app, cmd.Opts)

	cmdOpts := &CmdOpts{
		AppOpts:             app,
		CompiledReleaseOpts: cmd.Opts,
	}

	cmd.AnalyzersCmd.CmdOpts = cmdOpts
	cmd.ArtifactCmd.CmdOpts = cmdOpts
	cmd.OpsFileCmd.CmdOpts = cmdOpts
	cmd.UploadReleaseCmd.CmdOpts = cmdOpts
	cmd.ExportReleaseCmd.CmdOpts = cmdOpts
	cmd.DownloadCmd.CmdOpts = cmdOpts

	return cmd
}
