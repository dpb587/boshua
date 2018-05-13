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

	results, err := index.Filter(c.ReleaseOpts.Reference())
	if err != nil {
		return errors.Wrap(err, "filtering")
	}

	for _, result := range results {
		resultRef := result.Reference().(releaseversion.Reference)

		fmt.Printf("%s\t%s\n", resultRef.Name, resultRef.Version)
	}

	return nil
}
