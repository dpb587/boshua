package cli

import (
	"os"

	"github.com/dpb587/boshua/pivnetfile/analyzers/tilearchivefiles.v1/formatter"
	"github.com/dpb587/boshua/util/checksum/algorithm"
)

type Sha256sumCmd struct{}

func (c *Sha256sumCmd) Execute(_ []string) error {
	f := formatter.Shasum{Algorithm: algorithm.MustLookupName("sha256")}
	return f.Format(os.Stdout, os.Stdin)
}
