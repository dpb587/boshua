package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dpb587/boshua/api/v2/models"
	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/cli/client/args"
	"github.com/dpb587/boshua/manifest"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/util/metalinkutil"
)

type PatchManifestCmd struct {
	*CmdOpts `no-flag:"true"`

	Release     []string `long:"release" description:"Only check the release(s) matching this name (glob-friendly)"`
	SkipRelease []string `long:"skip-release" description:"Skip the release(s) matching this name (glob-friendly)"`

	LocalOS args.OS `long:"local-os" description:"Explicit local OS and version (used for bootstrap manifests)"`

	Parallel       int           `long:"parallel" description:"Maximum number of parallel operations"`
	RequestAndWait bool          `long:"request-and-wait" description:"Request and wait for compilations to finish"`
	WaitTimeout    time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`
}

func (c *PatchManifestCmd) Execute(_ []string) error {
	c.AppOpts.ConfigureLogger("patch-manifest")

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
		log.Fatalf("reading stdin: %v", err)
	}

	man, err := manifest.Parse(bytes, localStemcell)
	if err != nil {
		log.Fatalf("parsing manifest: %v", err)
	}

	apiclient := c.AppOpts.GetClient()

	for _, rel := range man.Requirements() {
		cs, err := checksum.CreateFromString(rel.Source.Sha1)
		if err != nil {
			log.Fatalf("parsing checksum: %v", err)
		}

		releaseVersionRef := models.ReleaseVersionRef{
			Name:     rel.Name,
			Version:  rel.Version,
			Checksum: cs,
		}
		osVersionRef := models.OSVersionRef{
			Name:    rel.Stemcell.OS,
			Version: rel.Stemcell.Version,
		}

		resInfo, err := apiclient.CompiledReleaseVersionInfo(models.CRVInfoRequest{
			Data: models.CRVInfoRequestData{
				ReleaseVersionRef: releaseVersionRef,
				OSVersionRef:      osVersionRef,
			},
		})
		if err != nil {
			log.Fatalf("finding compiled release: %v", err)
		} else if resInfo == nil {
			if !c.RequestAndWait {
				continue
			}

			priorStatus := "unknown"

			for {
				res, err := apiclient.CompiledReleaseVersionRequest(models.CRVRequestRequest{
					Data: models.CRVRequestRequestData{
						ReleaseVersionRef: releaseVersionRef,
						OSVersionRef:      osVersionRef,
					},
				})
				if err != nil {
					log.Fatalf("requesting compiled release: %v", err)
				} else if res == nil {
					fmt.Fprintf(os.Stderr, "[%s %s] unsupported compilation\n", rel.Stemcell.Slug(), rel.Slug())

					break
				}

				if res.Status != priorStatus {
					fmt.Fprintf(os.Stderr, "[%s %s] compilation status: %s\n", rel.Stemcell.Slug(), rel.Slug(), res.Status)
					priorStatus = res.Status
				}

				if res.Complete {
					break
				}

				time.Sleep(10 * time.Second)
			}

			if priorStatus == "unknown" {
				continue
			}

			resInfo, err = apiclient.CompiledReleaseVersionInfo(models.CRVInfoRequest{
				Data: models.CRVInfoRequestData{
					ReleaseVersionRef: releaseVersionRef,
					OSVersionRef:      osVersionRef,
				},
			})
			if err != nil {
				log.Fatalf("finding compiled release: %v", err)
			} else if resInfo == nil {
				log.Fatalf("finding compiled release: unable to verify request")
			}
		}

		rel.Compiled.Sha1 = metalinkutil.HashToChecksum(resInfo.Data.Artifact.Hashes[0]).String()
		rel.Compiled.URL = resInfo.Data.Artifact.URLs[0].URL

		err = man.UpdateRelease(rel)
		if err != nil {
			log.Fatalf("updating release: %v", err)
		}
	}

	bytes, err = man.Bytes()
	if err != nil {
		log.Fatalf("getting bytes: %v", err)
	}

	fmt.Printf("%s\n", bytes)

	return nil
}
