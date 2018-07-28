package cli

import (
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/boshua/metalink/metalinkutil"
)

type DownloadCmd struct {
	clicommon.DownloadCmd

	*CmdOpts `no-flag:"true"`

	Args DownloadCmdArgs `positional-args:"true" required:"true"`
}

type DownloadCmdArgs struct {
	Metalink string `positional-arg-name:"PATH" description:"Path to the metalink file"`
}

func (c *DownloadCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("artifact/download")

	return c.DownloadCmd.ExecuteArtifact(metalinkutil.NewStaticArtifactLoader(c.Args.Metalink))
}
