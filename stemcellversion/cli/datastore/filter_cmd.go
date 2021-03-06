package datastore

import (
	"fmt"

	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type FilterCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`
}

func (c *FilterCmd) Execute(_ []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "stemcell/datastore/filter"})

	index, err := c.Config.GetStemcellIndex(c.CmdOpts.DatastoreOpts.Datastore)
	if err != nil {
		return errors.Wrap(err, "loading stemcell index")
	}

	f, l := c.StemcellOpts.ArtifactParams()

	results, err := index.GetArtifacts(f, l)
	if err != nil {
		return errors.Wrap(err, "getting stemcell artifacts")
	}

	stemcellversion.Sort(results)

	for _, result := range results {
		fmt.Printf("%s\t%s\t%s\t%s\t%s\n", result.OS, result.Version, result.IaaS, result.Hypervisor, result.Flavor)
	}

	return nil
}
