package cli

import (
	cmdopts "github.com/dpb587/boshua/main/boshua/cmd/opts"
)

type Cmd struct {
	UseCompiledReleasesCmd UseCompiledReleasesCmd `command:"use-compiled-releases" description:"For patching a manifest to refer to compiled releases"`
}

func New(app *cmdopts.Opts) *Cmd {
	return &Cmd{}
}
