package cli

import (
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	"github.com/dpb587/boshua/compiledreleaseversion/cli/analysis"
	"github.com/dpb587/boshua/compiledreleaseversion/cli/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion/cli/opts"
)

type CmdOpts struct {
	AppOpts             *cmdopts.Opts `no-flag:"true"`
	CompiledReleaseOpts *opts.Opts
}

type Cmd struct {
	*opts.Opts

	AnalysisCmd  *analysis.Cmd  `command:"analysis" description:"For analyzing artifacts" subcommands-optional:"true"`
	DatastoreCmd *datastore.Cmd `command:"datastore" description:"For interacting with compiled release datastores"`

	AnalyzersCmd     AnalyzersCmd     `command:"analyzers" description:"For showing the supported analyzers"`
	ArtifactCmd      ArtifactCmd      `command:"artifact" description:"For showing the compiled release artifact"`
	OpsFileCmd       OpsFileCmd       `command:"ops-file" description:"For showing a deployment manifest ops file for the compiled release"`
	UploadReleaseCmd UploadReleaseCmd `command:"upload-release" description:"For uploading the compiled release to BOSH"`
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
		AppOpts:             app,
		CompiledReleaseOpts: cmd.Opts,
	}

	cmd.AnalyzersCmd.CmdOpts = cmdOpts
	cmd.ArtifactCmd.CmdOpts = cmdOpts
	cmd.OpsFileCmd.CmdOpts = cmdOpts
	cmd.UploadReleaseCmd.CmdOpts = cmdOpts

	return cmd
}
