package cli

import (
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/boshua/metalink/metalinkutil"
)

type UploadReleaseCmd struct {
	clicommon.UploadReleaseCmd

	*CmdOpts `no-flag:"true"`

	Name     string `long:"name" description:"Release name"`
	Version  string `long:"version" description:"Release version"`
	Stemcell string `long:"stemcell" description:"Compiled release stemcell (os/version format)"`

	Args UploadReleaseCmdArgs `positional-args:"true" required:"true"`
}

type UploadReleaseCmdArgs struct {
	Metalink string `positional-arg-name:"PATH" description:"Path to the metalink file"`
}

func (c *UploadReleaseCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("artifact/upload-release")

	return c.UploadReleaseCmd.ExecuteArtifact(metalinkutil.NewStaticArtifactLoader(c.Args.Metalink), clicommon.UploadReleaseOpts{
		Name:     c.Name,
		Version:  c.Version,
		Stemcell: c.Stemcell,
	})
}
