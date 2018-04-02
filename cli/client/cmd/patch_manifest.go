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

	man, err := manifest.Parse(bytes)
	if err != nil {
		log.Fatalf("parsing manifest: %v", err)
	}

	client := client.New(http.DefaultClient, c.Server)

	for _, rel := range man.Requirements() {
		releaseRef := models.ReleaseRef{
			Name:     rel.Name,
			Version:  rel.Version,
			Checksum: models.Checksum(fmt.Sprintf("sha1:%s", rel.Source.Sha1)),
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
		} else if resInfo == nil {
			if !c.RequestAndWait {
				continue
			}

			priorStatus := "unknown"

			for {
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
		}

		rel.Compiled.Sha1 = string(resInfo.Data.Tarball.Checksums[0])
		rel.Compiled.URL = resInfo.Data.Tarball.URL

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
