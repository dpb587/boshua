package cli

import (
	"github.com/dpb587/boshua/analysis/cli/formatter"
	cmdopts "github.com/dpb587/boshua/main/boshua/cmd/opts"
)

type CmdOpts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`
}

type Cmd struct {
	Formatter formatter.Cmd `command:"formatter" description:"For formatting the results of an analysis"`

	GenerateCmd GenerateCmd `command:"generate" description:"For generating an analysis"`
}
