package compiledrelease

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dpb587/bosh-compiled-releases/api/v2/client"
	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	"github.com/dpb587/bosh-compiled-releases/cli/client/args"
	"github.com/dpb587/metalink"
)

type Metalink struct {
	Server      string `long:"server" description:"Server address" default:"http://localhost:8080/" env:"CFBS_SERVER"`
	ServerToken string `long:"server-token" description:"Server authentication token" env:"CFBS_SERVER_TOKEN"`
	// CACert []string `long:"ca-cert" description:"Specific CA Certificate to trust"`

	RequestAndWait bool          `long:"request-and-wait" description:"Request and wait for compilations to finish"`
	WaitTimeout    time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`

	Quiet bool `long:"quiet" description:"Suppress informational output"`

	Args MetalinkArgs `positional-args:"true" optional:"true"`
}

type MetalinkArgs struct {
	Release         args.Release  `positional-arg-name:"RELEASE-NAME/RELEASE-VERSION"`
	Stemcell        args.Stemcell `positional-arg-name:"OS-NAME/STEMCELL-VERSION"`
	ReleaseChecksum args.Checksum `positional-arg-name:"RELEASE-CHECKSUM"`
}

func (c *Metalink) Execute(_ []string) error {
	apiclient := client.New(http.DefaultClient, c.Server)

	releaseRef := models.ReleaseRef{
		Name:     c.Args.Release.Name,
		Version:  c.Args.Release.Version,
		Checksum: models.Checksum(c.Args.ReleaseChecksum.String()),
	}
	stemcellRef := models.StemcellRef{
		OS:      c.Args.Stemcell.OS,
		Version: c.Args.Stemcell.Version,
	}

	var resInfo *models.CRVInfoResponse
	var err error

	if c.RequestAndWait {
		resInfo, err = client.RequestAndWait(apiclient, releaseRef, stemcellRef)
	} else {
		resInfo, err = apiclient.CompiledReleaseVersionInfo(models.CRVInfoRequest{
			Data: models.CRVInfoRequestData{
				Release:  releaseRef,
				Stemcell: stemcellRef,
			},
		})
	}

	if err != nil {
		log.Fatalf("requesting compiled version info: %v", err)
	} else if resInfo == nil {
		log.Fatalf("no compiled release available")
	}

	meta4 := createMetalink(resInfo)

	meta4Bytes, err := metalink.Marshal(meta4)
	if err != nil {
		log.Fatalf("marshalling response: %v", err)
	}

	fmt.Printf("%s\n", meta4Bytes)

	return nil
}
