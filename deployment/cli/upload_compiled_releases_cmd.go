package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"code.cloudfoundry.org/workpool"
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/artifact/cli/clicommon"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/dpb587/boshua/deployment/manifest"
	"github.com/dpb587/boshua/osversion"
	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type UploadCompiledReleasesCmd struct {
	setter.AppConfig `no-flag:"true"`

	Release     []string `long:"release" description:"Only check the release(s) matching this name (glob-friendly)"`
	SkipRelease []string `long:"skip-release" description:"Skip the release(s) matching this name (glob-friendly)"`

	Parallel int `long:"parallel" description:"Maximum number of parallel operations" default:"3"`

	clicommon.UploadReleaseCmd
}

func (c *UploadCompiledReleasesCmd) Execute(_ []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "deployment/upload-compiled-releases"})

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return errors.Wrap(err, "reading stdin")
	}

	man, err := manifest.Parse(bytes, osversion.Reference{})
	if err != nil {
		return errors.Wrap(err, "parsing manifest")
	}

	index, err := c.Config.GetReleaseCompilationIndex("default")
	if err != nil {
		return errors.Wrap(err, "loading index")
	}

	var parallelize []func()

	requirements := man.ReleaseRequirements()

	for relIdx := range requirements {
		rel := requirements[relIdx]

		parallelize = append(parallelize, func() {
			f := rel.FilterParams()

			parallelLog := func(msg string) {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%s [%s/%s %s/%s] %s\n", time.Now().Format("15:04:05"), rel.Stemcell.OS, rel.Stemcell.Version, rel.Name, rel.Version, msg))
			}

			result, err := compilationdatastore.GetCompilationArtifact(index, f)
			if err != nil {
				parallelLog(errors.Wrap(err, "skipped: error: getting compilation artifact").Error())

				return
			}

			err = c.UploadReleaseCmd.ExecuteArtifact(
				c.Config.GetDownloader,
				func() (artifact.Artifact, error) {
					return result, nil
				},
				clicommon.UploadReleaseOpts{},
			)
			if err != nil {
				log.Fatalf(fmt.Errorf("skipped: error: uploading release: %v", err).Error())

				return
			}

			parallelLog("uploaded compiled release")
		})
	}

	pool, err := workpool.NewThrottler(c.Parallel, parallelize)
	if err != nil {
		log.Fatalf(fmt.Errorf("parallelizing: %v", err).Error())
	}
	pool.Work()

	return nil
}
