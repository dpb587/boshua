package cli

import (
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/boshua/metalink/metalinkutil"
)

type UploadStemcellCmd struct {
	clicommon.UploadStemcellCmd

	*CmdOpts `no-flag:"true"`

	Name    string `long:"name" description:"Stemcell name"`
	Version string `long:"version" description:"Stemcell version"`

	Args UploadStemcellCmdArgs `positional-args:"true" required:"true"`
}

type UploadStemcellCmdArgs struct {
	Metalink string `positional-arg-name:"PATH" description:"Path to the metalink file"`
}

func (c *UploadStemcellCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("artifact/upload-stemcell")

	return c.UploadStemcellCmd.ExecuteArtifact(metalinkutil.NewStaticArtifactLoader(c.Args.Metalink), clicommon.UploadStemcellOpts{
		Name:    c.Name,
		Version: c.Version,
	})
}
