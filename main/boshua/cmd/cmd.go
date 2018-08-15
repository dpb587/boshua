package cmd

import (
	analysis "github.com/dpb587/boshua/analysis/cli"
	artifact "github.com/dpb587/boshua/artifact/cli"
	"github.com/dpb587/boshua/cli"
	"github.com/dpb587/boshua/cli/args"
	globalcmd "github.com/dpb587/boshua/cli/cmd"
	globalopts "github.com/dpb587/boshua/cli/opts"
	deployment "github.com/dpb587/boshua/deployment/cli"
	"github.com/dpb587/boshua/main/boshua/cmd/opts"
	releaseversion "github.com/dpb587/boshua/releaseversion/cli"
	server "github.com/dpb587/boshua/server/cli"
	stemcellversion "github.com/dpb587/boshua/stemcellversion/cli"
	"github.com/sirupsen/logrus"
)

type Cmd struct {
	App cli.App
	*opts.Opts

	AnalysisCmd   analysis.Cmd         `command:"analysis" description:"For analyzing artifacts"`
	ArtifactCmd   *artifact.Cmd        `command:"artifact" description:"For referencing artifacts"`
	ReleaseCmd    *releaseversion.Cmd  `command:"release" description:"For working with releases" subcommands-optional:"true"`
	DeploymentCmd *deployment.Cmd      `command:"deployment" description:"For working with deployments"`
	StemcellCmd   *stemcellversion.Cmd `command:"stemcell" description:"For working with stemcells" subcommands-optional:"true"`

	ServerCmd  server.Cmd           `command:"server" description:"For running an API server for remote access"`
	VersionCmd globalcmd.VersionCmd `command:"version" description:"For showing the version of this tool"`
}

func New(app cli.App) *Cmd {
	cmd := &Cmd{
		App: app,
		Opts: &opts.Opts{
			Opts: &globalopts.Opts{
				LogLevel: args.LogLevel(logrus.FatalLevel),
			},
		},
	}

	cmd.ArtifactCmd = artifact.New(cmd.Opts)
	cmd.ReleaseCmd = releaseversion.New(cmd.Opts)
	cmd.DeploymentCmd = deployment.New(cmd.Opts)
	cmd.StemcellCmd = stemcellversion.New(cmd.Opts)

	cmd.VersionCmd.App = cmd.App

	return cmd
}
