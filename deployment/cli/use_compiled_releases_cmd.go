package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"code.cloudfoundry.org/workpool"
	"github.com/dpb587/boshua/cli/args"
	"github.com/dpb587/boshua/deployment/manifest"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/osversion"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

type UseCompiledReleasesCmd struct {
	*CmdOpts `no-flag:"true"`

	Release     []string `long:"release" description:"Only check the release(s) matching this name (glob-friendly)"`
	SkipRelease []string `long:"skip-release" description:"Skip the release(s) matching this name (glob-friendly)"`

	LocalOS args.OS `long:"local-os" description:"Explicit local OS and version (used for bootstrap manifests)"`

	Parallel    int           `long:"parallel" description:"Maximum number of parallel operations" default:"3"`
	NoWait      bool          `long:"no-wait" description:"Do not request and wait for compilation if not already available"`
	WaitTimeout time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`
}

func (c *UseCompiledReleasesCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("deployment/patch-manifest")

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

	scheduler, err := c.AppOpts.GetScheduler()
	if err != nil {
		return errors.Wrap(err, "loading scheduler")
	}

	var parallelize []func()

	requirements := man.Requirements()

	for relIdx := range requirements {
		rel := requirements[relIdx]

		parallelize = append(parallelize, func() {
			f := compilationdatastore.FilterParams{
				Release: releaseversiondatastore.FilterParams{
					NameExpected:     true,
					Name:             rel.Name,
					VersionExpected:  true,
					Version:          rel.Version,
					ChecksumExpected: true,
					Checksum:         fmt.Sprintf("sha1:%s", rel.Source.Sha1),
				},
				OS: osversiondatastore.FilterParams{
					NameExpected:    true,
					Name:            rel.Stemcell.OS,
					VersionExpected: true,
					Version:         rel.Stemcell.Version,
				},
			}

			parallelLog := func(msg string) {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%s [%s/%s %s/%s] %s\n", time.Now().Format("15:04:05"), rel.Stemcell.OS, rel.Stemcell.Version, rel.Name, rel.Version, msg))
			}

			index, err := c.AppOpts.GetCompiledReleaseIndex("default")
			if err != nil {
				parallelLog(errors.Wrap(err, "loading index").Error())

				return
			}

			result, err := compilationdatastore.GetCompilationArtifact(index, f)
			if err == compilationdatastore.NoMatchErr {
				if c.NoWait {
					parallelLog(errors.New("none found").Error())

					return
				}

				releaseVersionIndex, err := c.AppOpts.GetReleaseIndex("default")
				if err != nil {
					parallelLog(errors.Wrap(err, "loading release index").Error())

					return
				}

				releaseVersion, err := releaseversiondatastore.GetArtifact(releaseVersionIndex, f.Release)
				if err != nil {
					parallelLog(errors.Wrap(err, "filtering release").Error())

					return
				}

				stemcellVersionIndex, err := c.AppOpts.GetStemcellIndex("default")
				if err != nil {
					parallelLog(errors.Wrap(err, "loading stemcell index").Error())

					return
				}

				stemcellVersion, err := stemcellversiondatastore.GetArtifact(stemcellVersionIndex, stemcellversiondatastore.FilterParams{
					OSExpected:      true,
					OS:              f.OS.Name,
					VersionExpected: true,
					Version:         f.OS.Version,
					// TODO dynamic
					IaaSExpected:   true,
					IaaS:           "aws",
					FlavorExpected: true,
					Flavor:         "light",
				})
				if err != nil {
					parallelLog(errors.Wrap(err, "filtering stemcell").Error())

					return
				}

				scheduledTask, err := scheduler.ScheduleCompilation(releaseVersion, stemcellVersion)
				if err != nil {
					parallelLog(errors.Wrap(err, "creating compilation").Error())

					return
				}

				status, err := schedulerpkg.WaitForScheduledTask(scheduledTask, func(status schedulerpkg.Status) {
					if c.AppOpts.Quiet {
						return
					}

					parallelLog(fmt.Sprintf("compilation is %s", status))
				})
				if err != nil {
					parallelLog(errors.Wrap(err, "checking task").Error())

					return
				} else if status != schedulerpkg.StatusSucceeded {
					parallelLog(fmt.Errorf("task did not succeed: %s", status).Error())

					return
				}

				result, err = compilationdatastore.GetCompilationArtifact(index, f)
				if err != nil {
					parallelLog(errors.Wrap(err, "finding finished compilation").Error())

					return
				}
			} else if err != nil {
				parallelLog(errors.Wrap(err, "filtering").Error())

				return
			}

			rel.Compiled.Sha1 = metalinkutil.HashToChecksum(result.Tarball.Hashes[0]).String()
			rel.Compiled.URL = result.Tarball.URLs[0].URL

			err = man.UpdateRelease(rel)
			if err != nil {
				log.Fatalf(fmt.Errorf("updating release: %v", err).Error())

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
