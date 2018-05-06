package cmd

import (
	"github.com/dpb587/boshua/cli/client/args"
	"github.com/dpb587/boshua/cli/client/cmd/analysis"
	"github.com/dpb587/boshua/cli/client/cmd/compiledrelease"
	"github.com/dpb587/boshua/cli/client/cmd/deployment"
	"github.com/dpb587/boshua/cli/client/cmd/opts"
	"github.com/dpb587/boshua/cli/client/cmd/release"
	"github.com/dpb587/boshua/cli/client/cmd/stemcell"
	"github.com/sirupsen/logrus"
)

type CmdOpts struct {
	AppOpts *opts.Opts
}

type Cmd struct {
	*opts.Opts

	AnalysisCmd        *analysis.Cmd        `command:"analysis" description:"For analyzing artifacts"`
	CompiledReleaseCmd *compiledrelease.Cmd `command:"compiled-release" description:"For working with compiled releases" subcommands-optional:"true"`
	ReleaseCmd         *release.Cmd         `command:"release" description:"For working with releases" subcommands-optional:"true"`
	DeploymentCmd      *deployment.Cmd      `command:"deployment" description:"For working with deployments"`
	StemcellCmd        *stemcell.Cmd        `command:"stemcell" description:"For working with stemcells" subcommands-optional:"true"`

	DownloadMetalinkCmd DownloadMetalinkCmd `command:"download-metalink" description:"Internal. Download resources in a metalink."`
}

func New() *Cmd {
	app := &Cmd{
		Opts: &opts.Opts{
			LogLevel: args.LogLevel(logrus.FatalLevel),
		},
	}

	app.AnalysisCmd = analysis.New(app.Opts)
	app.CompiledReleaseCmd = compiledrelease.New(app.Opts)
	app.ReleaseCmd = release.New(app.Opts)
	app.DeploymentCmd = deployment.New(app.Opts)
	app.StemcellCmd = stemcell.New(app.Opts)

	cmdOpts := &CmdOpts{
		AppOpts: app.Opts,
	}

	app.DownloadMetalinkCmd.CmdOpts = cmdOpts

	return app
}
