package compiledrelease

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dpb587/bosh-compiled-releases/api/v2/client"
	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	"github.com/dpb587/metalink"
)

type MetalinkCmd struct {
	*CmdOpts `no-flag:"true"`

	RequestAndWait bool          `long:"request-and-wait" description:"Request and wait for compilations to finish"`
	WaitTimeout    time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`

	Quiet bool `long:"quiet" description:"Suppress informational output"`
}

func (c *MetalinkCmd) Execute(_ []string) error {
	apiclient := client.New(http.DefaultClient, c.AppOpts.Server)

	releaseRef := models.ReleaseRef{
		Name:     c.CompiledReleaseOpts.Release.Name,
		Version:  c.CompiledReleaseOpts.Release.Version,
		Checksum: models.Checksum(c.CompiledReleaseOpts.ReleaseChecksum.String()),
	}
	stemcellRef := models.StemcellRef{
		OS:      c.CompiledReleaseOpts.Stemcell.OS,
		Version: c.CompiledReleaseOpts.Stemcell.Version,
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
