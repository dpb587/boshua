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
	"github.com/dpb587/boshua/task"
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

			index, err := c.AppOpts.GetCompiledReleaseIndex("default")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "loading index"))

				return
			}

			results, err := index.GetCompilationArtifacts(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "filtering"))

				return
			}

			if len(results) == 0 {
				if c.NoWait {
					fmt.Fprintf(os.Stderr, "%s\n", errors.New("none found"))

					return
				}

				releaseVersionIndex, err := c.AppOpts.GetReleaseIndex("default")
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "loading release index"))

					return
				}

				releaseVersions, err := releaseVersionIndex.GetArtifacts(f.Release)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "filtering release"))

					return
				}

				releaseVersion, err := releaseversiondatastore.RequireSingleResult(releaseVersions)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "filtering release"))

					return
				}

				stemcellVersionIndex, err := c.AppOpts.GetStemcellIndex("default")
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "loading stemcell index"))

					return
				}

				stemcellVersions, err := stemcellVersionIndex.GetArtifacts(stemcellversiondatastore.FilterParams{
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
					fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "filtering stemcell"))

					return
				}

				stemcellVersion, err := stemcellversiondatastore.RequireSingleResult(stemcellVersions)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "filtering stemcell"))

					return
				}

				scheduledTask, err := scheduler.ScheduleCompilation(releaseVersion, stemcellVersion)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "creating compilation"))

					return
				}

				status, err := task.WaitForScheduledTask(scheduledTask, func(status task.Status) {
					if c.AppOpts.Quiet {
						return
					}

					fmt.Fprintf(os.Stderr, "%s [%s/%s %s/%s] compilation is %s\n", time.Now().Format("15:04:05"), stemcellVersion.OS, stemcellVersion.Version, releaseVersion.Name, releaseVersion.Version, status)
				})
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "checking task"))

					return
				} else if status != task.StatusSucceeded {
					fmt.Fprintf(os.Stderr, "%s\n", fmt.Errorf("task did not succeed: %s", status))

					return
				}

				results, err = index.GetCompilationArtifacts(f)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "finding finished compilation"))

					return
				}
			}

			result, err := compilationdatastore.RequireSingleResult(results)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", errors.Wrap(err, "filtering"))

				return
			}

			rel.Compiled.Sha1 = metalinkutil.HashToChecksum(result.Tarball.Hashes[0]).String()
			rel.Compiled.URL = result.Tarball.URLs[0].URL

			err = man.UpdateRelease(rel)
			if err != nil {
				fmt.Fprintf(os.Stderr, "updating release: %v", err)

				return
			}
		})
	}

	pool, err := workpool.NewThrottler(c.Parallel, parallelize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parallelizing: %v", err)
	}
	pool.Work()

	bytes, err = man.Bytes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "getting bytes: %v", err)
	}

	fmt.Printf("%s\n", bytes)

	return nil
}
