package cmd

import (
	"github.com/dpb587/boshua/cli/client/cmd/analysis"
	"github.com/dpb587/boshua/cli/client/cmd/compiledrelease"
	"github.com/dpb587/boshua/cli/client/cmd/deployment"
	"github.com/dpb587/boshua/cli/client/cmd/opts"
	"github.com/dpb587/boshua/cli/client/cmd/release"
)

type CmdOpts struct {
	AppOpts *opts.Opts
}

type Cmd struct {
	*opts.Opts

	AnalysisCmd        *analysis.Cmd        `command:"analysis" description:"For analyzing artifacts"`
	CompiledReleaseCmd *compiledrelease.Cmd `command:"compiled-release" description:"For working with compiled releases"`
	ReleaseCmd         *release.Cmd         `command:"release" description:"For working with releases"`
	DeploymentCmd      *deployment.Cmd      `command:"deployment" description:"For working with deployments"`
}

func New() *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmd.AnalysisCmd = analysis.New(cmd.Opts)
	cmd.CompiledReleaseCmd = compiledrelease.New(cmd.Opts)
	cmd.ReleaseCmd = release.New(cmd.Opts)
	cmd.DeploymentCmd = deployment.New(cmd.Opts)

	return cmd
}
