package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dpb587/bosh-compiled-releases/api/v2/client"
	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	"github.com/dpb587/bosh-compiled-releases/manifest"
)

type PatchManifest struct {
	Server string `long:"server" description:"Server address" default:"http://localhost:8080/"`
	// CACert []string `long:"ca-cert" description:"Specific CA Certificate to trust"`

	RequestAndWait bool `long:"request-and-wait" description:"Request and wait for compilations to finish"`
	Quiet          bool `long:"quiet" description:"Suppress informational output"`
}

func (c *PatchManifest) Execute(_ []string) error {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("reading stdin: %v", err)
	}

	client := client.New(http.DefaultClient, c.Server)

	man, err := manifest.Parse(bytes)
	if err != nil {
		log.Fatalf("parsing manifest: %v", err)
	}

	for _, rel := range man.Requirements() {
		releaseRef := models.ReleaseRef{
			Name:    rel.Name,
			Version: rel.Version,
			Checksum: models.Checksum{
				Type:  "sha1",
				Value: rel.Source.Sha1,
			},
		}
		stemcellRef := models.StemcellRef{
			OS:      rel.Stemcell.OS,
			Version: rel.Stemcell.Version,
		}

		resInfo, err := client.CompiledReleaseVersionInfo(models.CRVInfoRequest{
			Data: models.CRVInfoRequestData{
				Release:  releaseRef,
				Stemcell: stemcellRef,
			},
		})
		if err != nil {
			log.Fatalf("finding compiled release: %v", err)
		} else if resInfo == nil || resInfo.Data.Status != "available" {
			if !c.RequestAndWait {
				continue
			}

			if resInfo == nil {
				res, err := client.CompiledReleaseVersionRequest(models.CRVRequestRequest{
					Data: models.CRVRequestRequestData{
						Release:  releaseRef,
						Stemcell: stemcellRef,
					},
				})
				if err != nil {
					log.Fatalf("requesting compiled release: %v", err)
				} else if res == nil {
					fmt.Fprintf(os.Stderr, "[%s %s] unsupported compilation\n", rel.Stemcell.Slug(), rel.Slug())

					continue
				}

				fmt.Fprintf(os.Stderr, "[%s %s] requested compiled release\n", rel.Stemcell.Slug(), rel.Slug())
			}

			fmt.Fprintf(os.Stderr, "[%s %s] waiting for compiled release\n", rel.Stemcell.Slug(), rel.Slug())

			for {
				time.Sleep(10 * time.Second)

				resInfo, err = client.CompiledReleaseVersionInfo(models.CRVInfoRequest{
					Data: models.CRVInfoRequestData{
						Release:  releaseRef,
						Stemcell: stemcellRef,
					},
				})
				if err != nil {
					log.Fatalf("finding compiled release: %v", err)
				} else if resInfo == nil {
					log.Fatalf("finding compiled release: unable to verify request")
				}

				if resInfo.Data.Status == "available" && resInfo.Data.Tarball.URL != "" {
					break
				}
			}
		}

		rel.Compiled.Sha1 = resInfo.Data.Tarball.Checksums[0].Value
		rel.Compiled.URL = resInfo.Data.Tarball.URL

		fmt.Printf("%#+v", rel)

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
