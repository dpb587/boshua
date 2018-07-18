package datastore

import (
	"fmt"

	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/pkg/errors"
)

type FilterCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *FilterCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("compiledrelease/datastore/filter")

	index, err := c.getDatastore()
	if err != nil {
		return errors.Wrap(err, "loading datastore")
	}

	results, err := index.Filter(c.CompiledReleaseOpts.FilterParams())
	if err != nil {
		return errors.Wrap(err, "filtering")
	}

	for _, result := range results {
		resultRef := result.Reference().(compilation.Reference)

		fmt.Printf(
			"%s\t%s\t%s\t%s\n",
			resultRef.ReleaseVersion.Name,
			resultRef.ReleaseVersion.Version,
			resultRef.OSVersion.Name,
			resultRef.OSVersion.Version,
		)
	}

	return nil
}
