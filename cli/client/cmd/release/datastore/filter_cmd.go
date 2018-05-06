package datastore

import (
	"fmt"
)

type FilterCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *FilterCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/datastore/filter")

	releaseIndex, err := c.getDatastore()
	if err != nil {
		return fmt.Errorf("loading datastore: %v", err)
	}

	results, err := releaseIndex.Filter(c.ReleaseOpts.Reference())
	if err != nil {
		return fmt.Errorf("filtering: %v", err)
	}

	for _, result := range results {
		fmt.Printf("%s\t%s\n", result.Name, result.Version)
	}

	return nil
}
