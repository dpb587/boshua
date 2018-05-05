package analyzer

import (
	"os"

	"github.com/dpb587/boshua/analysis/analyzer/stemcellimagefiles.v1/formatter"
	"github.com/dpb587/boshua/checksum/algorithm"
)

type Sha512sumCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *Sha512sumCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("analysis/stemcellimagefiles.v1/sha512sum")

	f := formatter.Shasum{Algorithm: algorithm.MustLookupName("sha512")}
	return f.Format(os.Stdout, os.Stdin)
}
