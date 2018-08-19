package datastore

import (
	"fmt"
	"path/filepath"

	"github.com/dpb587/boshua/config"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type StoreCmd struct {
	setter.AppConfig `no-flag:"true"`
	*CmdOpts         `no-flag:"true"`

	Version string `long:"version" description:"A specific version to use" default:"0.0.0"`

	Args StoreCmdArgs `positional-args:"true" required:"true"`
}

type StoreCmdArgs struct {
	Artifact string `positional-arg-name:"PATH" description:"Path to the artifact"`
}

func (c *StoreCmd) Execute(_ []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "compiledrelease/datastore/store"})

	index, err := c.Config.GetReleaseCompilationIndex(c.CmdOpts.DatastoreOpts.Datastore)
	if err != nil {
		return errors.Wrap(err, "loading datastore")
	}

	releaseVersion, err := c.CompiledReleaseOpts.ReleaseOpts.Artifact(c.AppConfig.Config)
	if err != nil {
		return errors.Wrap(err, "finding release")
	}

	osVersionIndex, err := c.Config.GetOSIndex(config.DefaultName)
	if err != nil {
		return errors.Wrap(err, "loading os index")
	}

	osVersion, err := osVersionIndex.Find(osversion.Reference{Name: c.CompiledReleaseOpts.OS.Name, Version: c.CompiledReleaseOpts.OS.Version})
	if err != nil {
		return errors.Wrap(err, "finding os")
	}

	path, err := filepath.Abs(c.Args.Artifact)
	if err != nil {
		return errors.Wrap(err, "expanding artifact path")
	}

	meta4, err := metalinkutil.CreateFromFiles(fmt.Sprintf("file://%s", path))
	if err != nil {
		return errors.Wrap(err, "building metalink")
	}

	return index.StoreCompilationArtifact(compilation.Artifact{
		Release: releaseVersion.Reference().(releaseversion.Reference),
		OS:      osVersion.Reference().(osversion.Reference),
		Tarball: meta4.Files[0],
	})
}
