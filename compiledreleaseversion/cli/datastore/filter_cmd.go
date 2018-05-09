package datastore

import (
	"fmt"
)

type FilterCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *FilterCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("compiledrelease/datastore/filter")

	index, err := c.getDatastore()
	if err != nil {
		return fmt.Errorf("loading datastore: %v", err)
	}

	results, err := index.Filter(c.CompiledReleaseOpts.Reference())
	if err != nil {
		return fmt.Errorf("filtering: %v", err)
	}

	for _, result := range results {
		fmt.Printf(
			"%s\t%s\t%s\t%s\n",
			result.ReleaseVersion.Name,
			result.ReleaseVersion.Version,
			result.OSVersion.Name,
			result.OSVersion.Version,
		)
	}

	return nil
}
