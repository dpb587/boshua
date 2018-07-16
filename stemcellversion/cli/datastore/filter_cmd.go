package datastore

import (
	"fmt"

	"github.com/dpb587/boshua/stemcellversion"
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
		resultRef := result.Reference().(stemcellversion.Reference)

		fmt.Printf("%s\t%s\t%s\t%s\n", resultRef.IaaS, resultRef.Hypervisor, resultRef.OS, resultRef.Version)
	}

	return nil
}
