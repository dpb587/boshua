package datastore

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

type LabelsCmd struct {
	*CmdOpts `no-flag:"true"`
}

func (c *LabelsCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("release/datastore/filter")

	index, err := c.getDatastore()
	if err != nil {
		return errors.Wrap(err, "loading datastore")
	}

	results, err := index.Labels()
	if err != nil {
		return errors.Wrap(err, "getting labels")
	}

	sort.Strings(results)

	for _, result := range results {
		fmt.Printf("%s\n", result)
	}

	return nil
}
