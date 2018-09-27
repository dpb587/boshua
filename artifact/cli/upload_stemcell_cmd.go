package cli

import (
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/config/provider/setter"
)

type UploadStemcellCmd struct {
	setter.AppConfig`no-flag:"true"`
	clicommon.UploadStemcellCmd

	Name    string `long:"name" description:"Stemcell name"`
	Version string `long:"version" description:"Stemcell version"`

	Args UploadStemcellCmdArgs `positional-args:"true" required:"true"`
}

type UploadStemcellCmdArgs struct {
	Metalink string `positional-arg-name:"PATH" description:"Path to the metalink file"`
}

func (c *UploadStemcellCmd) Execute(extra []string) error {
	return c.UploadStemcellCmd.ExecuteArtifact(
		c.Config.GetDownloader,
		metalinkutil.NewStaticArtifactLoader(c.Args.Metalink),
		clicommon.UploadStemcellOpts{
			Name:      c.Name,
			Version:   c.Version,
			ExtraArgs: extra,
		},
	)
}
