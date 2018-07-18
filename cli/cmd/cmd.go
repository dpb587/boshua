package cmd

import (
	analysis "github.com/dpb587/boshua/analysis/cli"
	"github.com/dpb587/boshua/cli/args"
	"github.com/dpb587/boshua/cli/cmd/opts"
	// deployment "github.com/dpb587/boshua/deployment/cli"
	releaseversion "github.com/dpb587/boshua/releaseversion/cli"
	server "github.com/dpb587/boshua/server/cli"
	stemcellversion "github.com/dpb587/boshua/stemcellversion/cli"
	"github.com/sirupsen/logrus"
)

type CmdOpts struct {
	AppOpts *opts.Opts
}

type Cmd struct {
	*opts.Opts

	AnalysisCmd *analysis.Cmd       `command:"analysis" description:"For analyzing artifacts"`
	ReleaseCmd  *releaseversion.Cmd `command:"release" description:"For working with releases" subcommands-optional:"true"`
	// DeploymentCmd      *deployment.Cmd             `command:"deployment" description:"For working with deployments"`
	StemcellCmd *stemcellversion.Cmd `command:"stemcell" description:"For working with stemcells" subcommands-optional:"true"`

	DownloadMetalinkCmd DownloadMetalinkCmd `command:"download-metalink" description:"For downloading assets in a metalink"`
	ServerCmd           server.Cmd          `command:"server" description:"For running an API server for remote access"`
}

func New() *Cmd {
	app := &Cmd{
		Opts: &opts.Opts{
			LogLevel: args.LogLevel(logrus.FatalLevel),
		},
	}

	app.AnalysisCmd = analysis.New(app.Opts)
	app.ReleaseCmd = releaseversion.New(app.Opts)
	// app.DeploymentCmd = deployment.New(app.Opts)
	app.StemcellCmd = stemcellversion.New(app.Opts)

	cmdOpts := &CmdOpts{
		AppOpts: app.Opts,
	}

	app.DownloadMetalinkCmd.CmdOpts = cmdOpts
	app.ServerCmd.AppOpts = app.Opts

	return app
}
