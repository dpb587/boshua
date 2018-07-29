package cli

import (
	"os"

	"github.com/dpb587/boshua/stemcellversion/analyzers/stemcellimagefiles.v1/formatter"
	"github.com/dpb587/boshua/util/checksum/algorithm"
)

type Sha512sumCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *Sha512sumCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("analysis/formatter/stemcellimagefiles.v1/sha512sum")

	f := formatter.Shasum{Algorithm: algorithm.MustLookupName("sha512")}
	return f.Format(os.Stdout, os.Stdin)
}
