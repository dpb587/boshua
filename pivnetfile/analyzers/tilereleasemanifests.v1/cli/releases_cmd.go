package cli

import (
	"os"

	"github.com/dpb587/boshua/pivnetfile/analyzers/tilereleasemanifests.v1/formatter"
)

type ReleasesCmd struct{}

func (c *ReleasesCmd) Execute(_ []string) error {
	f := formatter.Releases{}
	return f.Format(os.Stdout, os.Stdin)
}
