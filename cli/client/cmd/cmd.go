package cmd

import (
	"github.com/dpb587/bosh-compiled-releases/cli/client/cmd/analysis"
	"github.com/dpb587/bosh-compiled-releases/cli/client/cmd/compiledrelease"
	"github.com/dpb587/bosh-compiled-releases/cli/client/cmd/opts"
)

type CmdOpts struct {
	AppOpts *opts.Opts
}

type Cmd struct {
	*opts.Opts

	AnalysisCmd        *analysis.Cmd        `command:"analysis" description:"For analyzing artifacts"`
	CompiledReleaseCmd *compiledrelease.Cmd `command:"compiled-release" description:"For working with compiled releases"`
	PatchManifestCmd   PatchManifestCmd     `command:"patch-manifest" description:"For patching a manifest to use compiled releases"`
}

func New() *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmd.AnalysisCmd = analysis.New(cmd.Opts)
	cmd.CompiledReleaseCmd = compiledrelease.New(cmd.Opts)

	cmdOpts := &CmdOpts{
		AppOpts: cmd.Opts,
	}

	cmd.PatchManifestCmd.CmdOpts = cmdOpts

	return cmd
}
