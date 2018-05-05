package analyzer

import (
	"os"

	"github.com/dpb587/boshua/analysis/analyzer/stemcellimagefiles.v1/formatter"
	"github.com/dpb587/boshua/checksum/algorithm"
)

type Sha256sumCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *Sha256sumCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("analysis/stemcellimagefiles.v1/sha256sum")

	f := formatter.Shasum{Algorithm: algorithm.MustLookupName("sha256")}
	return f.Format(os.Stdout, os.Stdin)
}
