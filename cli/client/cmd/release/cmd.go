package release

import (
	cmdopts "github.com/dpb587/boshua/cli/client/cmd/opts"
	"github.com/dpb587/boshua/cli/client/cmd/release/analysis"
	"github.com/dpb587/boshua/cli/client/cmd/release/opts"
)

type CmdOpts struct {
	AppOpts     *cmdopts.Opts `no-flag:"true"`
	ReleaseOpts *opts.Opts
}

type Cmd struct {
	*opts.Opts

	AnalysisCmd *analysis.Cmd `command:"analysis" description:"For analyzing the release artifact"`

	ArtifactCmd      ArtifactCmd      `command:"artifact" description:"For showing the release artifact"`
	UploadReleaseCmd UploadReleaseCmd `command:"upload-release" description:"For uploading the release to BOSH"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.ArtifactCmd.Execute(extra)
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmd.AnalysisCmd = analysis.New(app, cmd.Opts)

	cmdOpts := &CmdOpts{
		AppOpts:     app,
		ReleaseOpts: cmd.Opts,
	}

	cmd.ArtifactCmd.CmdOpts = cmdOpts
	cmd.UploadReleaseCmd.CmdOpts = cmdOpts

	return cmd
}
