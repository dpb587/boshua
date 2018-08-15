package datastore

import (
	"fmt"
	"sort"

	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type LabelsCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`
}

func (c *LabelsCmd) Execute(_ []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "release/datastore/labels"})

	index, err := c.Config.GetReleaseIndex("default")
	if err != nil {
		return errors.Wrap(err, "loading datastore")
	}

	results, err := index.GetLabels()
	if err != nil {
		return errors.Wrap(err, "getting labels")
	}

	sort.Strings(results)

	for _, result := range results {
		fmt.Printf("%s\n", result)
	}

	return nil
}
