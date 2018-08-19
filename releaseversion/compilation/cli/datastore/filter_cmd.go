package datastore

import (
	"fmt"

	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type FilterCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`
}

func (c *FilterCmd) Execute(_ []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "compiledrelease/datastore/filter"})

	index, err := c.Config.GetReleaseCompilationIndex(c.CmdOpts.DatastoreOpts.Datastore)
	if err != nil {
		return errors.Wrap(err, "loading datastore")
	}

	results, err := index.GetCompilationArtifacts(c.CompiledReleaseOpts.FilterParams())
	if err != nil {
		return errors.Wrap(err, "filtering")
	}

	compilation.Sort(results)

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
