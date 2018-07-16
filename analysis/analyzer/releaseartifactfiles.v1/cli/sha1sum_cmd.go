package cli

import (
	"os"

	"github.com/dpb587/boshua/analysis/analyzer/releaseartifactfiles.v1/formatter"
	"github.com/dpb587/boshua/util/checksum/algorithm"
)

type Sha1sumCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *Sha1sumCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("analysis/releaseartifactfiles.v1/sha1sum")

	f := formatter.Shasum{Algorithm: algorithm.MustLookupName("sha1")}
	return f.Format(os.Stdout, os.Stdin)
}