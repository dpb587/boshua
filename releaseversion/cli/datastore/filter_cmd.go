package datastore

import (
	"fmt"

	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type FilterCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`
}

func (c *FilterCmd) Execute(_ []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "release/datastore/filter"})

	index, err := c.Config.GetReleaseIndex(c.CmdOpts.DatastoreOpts.Datastore)
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
