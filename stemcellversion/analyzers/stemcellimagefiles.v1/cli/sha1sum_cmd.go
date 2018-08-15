package cli

import (
	"os"

	"github.com/dpb587/boshua/stemcellversion/analyzers/stemcellimagefiles.v1/formatter"
	"github.com/dpb587/boshua/util/checksum/algorithm"
)

type Sha1sumCmd struct{}

func (c *Sha1sumCmd) Execute(_ []string) error {
	f := formatter.Shasum{Algorithm: algorithm.MustLookupName("sha1")}
	return f.Format(os.Stdout, os.Stdin)
}
