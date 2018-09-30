package datastore

import (
	"fmt"

	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/dpb587/boshua/pivnetfile"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type FilterCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`
}

func (c *FilterCmd) Execute(_ []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "pivnetfile/datastore/filter"})

	index, err := c.Config.GetPivnetFileIndex(c.CmdOpts.DatastoreOpts.Datastore)
	if err != nil {
		return errors.Wrap(err, "loading pivnet file index")
	}

	results, err := index.GetArtifacts(c.PivnetFileOpts.FilterParams())
	if err != nil {
		return errors.Wrap(err, "getting artifacts")
	}

	pivnetfile.Sort(results)

	for _, result := range results {
		fmt.Printf("%s\t%d\t%s\t%d\t%s\n", result.ProductSlug, result.ReleaseID, result.ReleaseVersion, result.FileID, result.File.Name)
	}

	return nil
}
