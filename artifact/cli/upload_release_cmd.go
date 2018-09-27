package cli

import (
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/config/provider/setter"
)

type UploadReleaseCmd struct {
	setter.AppConfig`no-flag:"true"`
	clicommon.UploadReleaseCmd

	Name     string `long:"name" description:"Release name"`
	Version  string `long:"version" description:"Release version"`
	Stemcell string `long:"stemcell" description:"Compiled release stemcell (os/version format)"`

	Args UploadReleaseCmdArgs `positional-args:"true" required:"true"`
}

type UploadReleaseCmdArgs struct {
	Metalink string `positional-arg-name:"PATH" description:"Path to the metalink file"`
}

func (c *UploadReleaseCmd) Execute(extra []string) error {
	return c.UploadReleaseCmd.ExecuteArtifact(
		c.Config.GetDownloader,
		metalinkutil.NewStaticArtifactLoader(c.Args.Metalink),
		clicommon.UploadReleaseOpts{
			Name:      c.Name,
			Version:   c.Version,
			Stemcell:  c.Stemcell,
			ExtraArgs: extra,
		},
	)
}
