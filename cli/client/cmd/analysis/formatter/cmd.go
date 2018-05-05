package formatter

import (
	stemcellimagefilesv1 "github.com/dpb587/boshua/cli/client/cmd/analysis/formatter/stemcellimagefiles.v1"
	stemcellpackagesv1 "github.com/dpb587/boshua/cli/client/cmd/analysis/formatter/stemcellpackages.v1"
	cmdopts "github.com/dpb587/boshua/cli/client/cmd/opts"
)

type Cmd struct {
	StemcellimagefilesV1 *stemcellimagefilesv1.Cmd `command:"stemcellimagefiles.v1" subcommands-optional:"true"`
	StemcellpackagesV1   *stemcellpackagesv1.Cmd   `command:"stemcellpackages.v1" subcommands-optional:"true"`
}

type CmdOpts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`
}

func New(app *cmdopts.Opts) *Cmd {
	cmd := &Cmd{
		StemcellimagefilesV1: stemcellimagefilesv1.New(app),
		StemcellpackagesV1:   stemcellpackagesv1.New(app),
	}

	return cmd
}
