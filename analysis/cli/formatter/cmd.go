package formatter

import (
	stemcellimagefilesv1 "github.com/dpb587/boshua/analysis/analyzer/stemcellimagefiles.v1/cli"
	stemcellmanifestv1 "github.com/dpb587/boshua/analysis/analyzer/stemcellmanifest.v1/cli"
	stemcellpackagesv1 "github.com/dpb587/boshua/analysis/analyzer/stemcellpackages.v1/cli"
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
)

type Cmd struct {
	StemcellimagefilesV1 *stemcellimagefilesv1.Cmd `command:"stemcellimagefiles.v1" subcommands-optional:"true"`
	StemcellpackagesV1   *stemcellpackagesv1.Cmd   `command:"stemcellpackages.v1" subcommands-optional:"true"`
	StemcellmanifestV1   *stemcellmanifestv1.Cmd   `command:"stemcellmanifest.v1" subcommands-optional:"true"`
}

type CmdOpts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{
		StemcellimagefilesV1: stemcellimagefilesv1.New(app),
		StemcellpackagesV1:   stemcellpackagesv1.New(app),
		StemcellmanifestV1:   stemcellmanifestv1.New(app),
	}

	return cmd
}
