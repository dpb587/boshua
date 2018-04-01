package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dpb587/bosh-compiled-releases/api/v2/client"
	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	"github.com/dpb587/metalink"
)

type Metalink struct {
	Server string `long:"server" description:"Server address" default:"http://localhost:8080/"`
	// CACert []string `long:"ca-cert" description:"Specific CA Certificate to trust"`

	RequestAndWait bool `long:"request-and-wait" description:"Request and wait for compilations to finish"`
	Quiet          bool `long:"quiet" description:"Suppress informational output"`

	Args MetalinkArgs `positional-args:"true" optional:"true"`
}

type MetalinkArgs struct {
	Release         string `positional-arg-name:"RELEASE-NAME/RELEASE-VERSION"`
	Stemcell        string `positional-arg-name:"OS-NAME/STEMCELL-VERSION"`
	ReleaseChecksum string `positional-arg-name:"RELEASE-CHECKSUM"`
}

func (c *Metalink) Execute(_ []string) error {
	releaseSplit := strings.SplitN(c.Args.Release, "/", 2)
	if len(releaseSplit) != 2 {
		log.Fatalf("expected name/version-formatted value: %s", c.Args.Stemcell)
	}

	stemcellSplit := strings.SplitN(c.Args.Stemcell, "/", 2)
	if len(stemcellSplit) != 2 {
		log.Fatalf("expected os/version-formatted value: %s", c.Args.Stemcell)
	}

	client := client.New(http.DefaultClient, c.Server)

	releaseRef := models.ReleaseRef{
		Name:     releaseSplit[0],
		Version:  releaseSplit[1],
		Checksum: models.Checksum(fmt.Sprintf("sha1:%s", c.Args.ReleaseChecksum)),
	}
	stemcellRef := models.StemcellRef{
		OS:      stemcellSplit[0],
		Version: stemcellSplit[1],
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
			log.Fatalf("compiled release is not available")
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
				log.Fatalf("unsupported compilation\n")
			}

			fmt.Fprintf(os.Stderr, "requested compiled release\n")
		}

		fmt.Fprintf(os.Stderr, "waiting for compiled release\n")

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

	meta4 := metalink.Metalink{
		Files: []metalink.File{
			{
				Name:    fmt.Sprintf("%s-%s-on-%s-version-%s.tgz", releaseRef.Name, releaseRef.Version, stemcellRef.OS, stemcellRef.Version),
				Version: releaseRef.Version,
				URLs: []metalink.URL{
					{
						URL: resInfo.Data.Tarball.URL,
					},
				},
			},
		},
		Generator: "bosh-compiled-releases/0.0.0",
	}

	if resInfo.Data.Tarball.Size != nil {
		meta4.Files[0].Size = *resInfo.Data.Tarball.Size
	}

	if resInfo.Data.Tarball.Published != nil {
		meta4.Published = resInfo.Data.Tarball.Published
	}

	for _, checksum := range resInfo.Data.Tarball.Checksums {
		var csType string

		switch checksum.Algorithm() {
		case "md5":
			csType = "md5"
		case "sha1":
			csType = "sha-1"
		case "sha256":
			csType = "sha-256"
		case "sha512":
			csType = "sha-512"
		default:
			continue
		}

		meta4.Files[0].Hashes = append(meta4.Files[0].Hashes, metalink.Hash{
			Type: csType,
			Hash: checksum.Data(),
		})
	}

	meta4Bytes, err := metalink.Marshal(meta4)
	if err != nil {
		log.Fatalf("marshalling response: %v", err)
	}

	fmt.Printf("%s\n", meta4Bytes)

	return nil
}
