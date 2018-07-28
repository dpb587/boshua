package datastore

import (
	"fmt"

	"github.com/dpb587/boshua/releaseversion"
	"github.com/pkg/errors"
)

type FilterCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *FilterCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/datastore/filter")

	index, err := c.getDatastore()
	if err != nil {
		return errors.Wrap(err, "loading datastore")
	}

	results, err := index.GetArtifacts(c.ReleaseOpts.FilterParams())
	if err != nil {
		return errors.Wrap(err, "filtering")
	}

	releaseversion.Sort(results)

	for _, result := range results {
		fmt.Printf("%s\t%s\n", result.Name, result.Version)
	}

	return nil
}
