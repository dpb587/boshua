package cli

import (
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/boshua/metalink/metalinkutil"
)

type DownloadCmd struct {
	clicommon.DownloadCmd `no-flag:"true"`

	*CmdOpts `no-flag:"true"`

	Args DownloadCmdArgs `positional-args:"true" required:"true"`
}

type DownloadCmdArgs struct {
	Metalink  string  `positional-arg-name:"PATH" description:"Path to the metalink file"`
	TargetDir *string `positional-arg-name:"TARGET-DIR" description:"Directory to download files (default: .)"`
}

func (c *DownloadCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("artifact/download")

	// cheat and inject targetdir
	c.DownloadCmd.Args.TargetDir = c.Args.TargetDir

	return c.DownloadCmd.ExecuteArtifact(metalinkutil.NewStaticArtifactLoader(c.Args.Metalink))
}
