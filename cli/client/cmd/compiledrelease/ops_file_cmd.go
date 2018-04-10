package compiledrelease

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dpb587/bosh-compiled-releases/api/v2/client"
	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	yaml "gopkg.in/yaml.v2"
)

type OpsFileCmd struct {
	*CmdOpts `no-flag:"true"`

	RequestAndWait bool          `long:"request-and-wait" description:"Request and wait for compilations to finish"`
	WaitTimeout    time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`

	Quiet bool `long:"quiet" description:"Suppress informational output"`
}

func (c *OpsFileCmd) Execute(_ []string) error {
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

	opsBytes, err := yaml.Marshal([]map[string]interface{}{
		{
			"path": fmt.Sprintf("/releases/name=%s?", resInfo.Data.Release.Name),
			"type": "replace",
			"value": map[string]interface{}{
				"name":    resInfo.Data.Release.Name,
				"version": resInfo.Data.Release.Version,
				"sha1":    strings.TrimPrefix(string(resInfo.Data.Tarball.Checksums[0]), "sha1:"), // TODO .Preferred()
				"url":     resInfo.Data.Tarball.URL,
				"stemcell": map[string]string{
					"os":      resInfo.Data.Stemcell.OS,
					"version": resInfo.Data.Stemcell.Version,
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("marshaling ops: %v", err)
	}

	fmt.Printf("%s\n", opsBytes)

	return nil
}
