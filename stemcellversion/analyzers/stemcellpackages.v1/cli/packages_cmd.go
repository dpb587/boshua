package cli

import (
	"os"

	"github.com/dpb587/boshua/stemcellversion/analyzers/stemcellpackages.v1/formatter"
)

type PackagesCmd struct{}

func (c *PackagesCmd) Execute(_ []string) error {
	f := formatter.Packages{}
	return f.Format(os.Stdout, os.Stdin)
}
