package cmd

import (
	"github.com/dpb587/boshua/cli/client/cmd/analysis"
	"github.com/dpb587/boshua/cli/client/cmd/compiledrelease"
	"github.com/dpb587/boshua/cli/client/cmd/deployment"
	"github.com/dpb587/boshua/cli/client/cmd/opts"
	"github.com/dpb587/boshua/cli/client/cmd/release"
	"github.com/dpb587/boshua/cli/client/cmd/stemcell"
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
}

func New() *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmd.AnalysisCmd = analysis.New(cmd.Opts)
	cmd.CompiledReleaseCmd = compiledrelease.New(cmd.Opts)
	cmd.ReleaseCmd = release.New(cmd.Opts)
	cmd.DeploymentCmd = deployment.New(cmd.Opts)
	cmd.StemcellCmd = stemcell.New(cmd.Opts)

	return cmd
}
