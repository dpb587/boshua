package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"code.cloudfoundry.org/workpool"
	"github.com/dpb587/boshua/cli/args"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/dpb587/boshua/deployment/manifest"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/osversion"
	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type UseCompiledReleasesCmd struct {
	setter.AppConfig `no-flag:"true"`

	Release     []string `long:"release" description:"Only check the release(s) matching this name (glob-friendly)"`
	SkipRelease []string `long:"skip-release" description:"Skip the release(s) matching this name (glob-friendly)"`

	LocalOS args.OS `long:"local-os" description:"Explicit local OS and version (used for bootstrap manifests)"`

	Parallel int `long:"parallel" description:"Maximum number of parallel operations" default:"3"`
}

func (c *UseCompiledReleasesCmd) Execute(_ []string) error {
	c.Config.AppendLoggerFields(logrus.Fields{"cli.command": "deployment/use-compiled-releases"})

	localStemcell := osversion.Reference{
		Name:    c.LocalOS.Name,
		Version: c.LocalOS.Version,
	}

	if localStemcell.Version == "" {
		bytes, err := ioutil.ReadFile("/var/vcap/bosh/etc/stemcell_version")
		if err != nil {
			if _, ok := err.(*os.PathError); !ok {
				log.Fatalf("reading stemcell_version")
			}
		}

		localStemcell.Version = string(bytes)
	}

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return errors.Wrap(err, "reading stdin")
	}

	man, err := manifest.Parse(bytes, localStemcell)
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

			rel.Compiled.Sha1 = metalinkutil.HashToChecksum(result.Tarball.Hashes[0]).String()
			rel.Compiled.URL = result.Tarball.URLs[0].URL

			err = man.UpdateRelease(rel)
			if err != nil {
				log.Fatalf(fmt.Errorf("skipped: error: updating release: %v", err).Error())

				return
			}

			parallelLog("added compiled release")
		})
	}

	pool, err := workpool.NewThrottler(c.Parallel, parallelize)
	if err != nil {
		log.Fatalf(fmt.Errorf("parallelizing: %v", err).Error())
	}
	pool.Work()

	bytes, err = man.Bytes()
	if err != nil {
		log.Fatalf(fmt.Errorf("getting bytes: %v", err).Error())
	}

	fmt.Printf("%s\n", bytes)

	return nil
}
