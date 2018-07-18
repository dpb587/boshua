package datastore

import (
	"fmt"

	"github.com/pkg/errors"
)

type FilterCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *FilterCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("stemcell/datastore/filter")

	index, err := c.getDatastore()
	if err != nil {
		return errors.Wrap(err, "loading datastore")
	}

	results, err := index.Filter(c.StemcellOpts.FilterParams())
	if err != nil {
		return errors.Wrap(err, "filtering")
	}

	for _, result := range results {
		fmt.Printf("%s\t%s\t%s\t%s\n", result.IaaS, result.Hypervisor, result.OS, result.Version)
	}

	return nil
}
