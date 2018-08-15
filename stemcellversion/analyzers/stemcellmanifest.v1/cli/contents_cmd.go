package cli

import (
	"os"

	"github.com/dpb587/boshua/stemcellversion/analyzers/stemcellmanifest.v1/formatter"
)

type ContentsCmd struct{}

func (c *ContentsCmd) Execute(_ []string) error {
	f := formatter.Contents{}
	return f.Format(os.Stdout, os.Stdin)
}
