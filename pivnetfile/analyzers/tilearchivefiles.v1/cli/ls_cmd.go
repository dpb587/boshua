package cli

import (
	"os"

	"github.com/dpb587/boshua/pivnetfile/analyzers/tilearchivefiles.v1/formatter"
)

type LsCmd struct{}

func (c *LsCmd) Execute(_ []string) error {
	f := formatter.Ls{}
	return f.Format(os.Stdout, os.Stdin)
}
