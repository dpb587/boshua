package cli

import (
	"io/ioutil"

	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
)

type UploadStemcellCmd struct {
	clicommon.UploadStemcellCmd

	*CmdOpts `no-flag:"true"`

	Args UploadStemcellCmdArgs `positional-args:"true" required:"true"`
}

type UploadStemcellCmdArgs struct {
	Metalink string `positional-arg-name:"PATH" description:"Path to the metalink file"`
}

func (c *UploadStemcellCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("artifact/upload-stemcell")

	return c.UploadStemcellCmd.ExecuteArtifact(func() (artifact.Artifact, error) {
		metalinkBytes, err := ioutil.ReadFile(c.Args.Metalink)
		if err != nil {
			return nil, errors.Wrap(err, "reading metalink")
		}

		var meta4 metalink.Metalink

		err = metalink.Unmarshal(metalinkBytes, &meta4)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshaling metalink")
		}

		return artifact.StaticArtifact{
			StaticMetalinkFile: meta4.Files[0],
		}, nil
	})
}