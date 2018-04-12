package opts

import (
	"time"

	"github.com/dpb587/bosh-compiled-releases/api/v2/client"
	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	"github.com/dpb587/bosh-compiled-releases/cli/client/args"
)

type Opts struct {
	Release         args.Release  `long:"release" description:"The release name and version"`
	ReleaseChecksum args.Checksum `long:"release-checksum" description:"The release checksum"`
	Stemcell        args.Stemcell `long:"stemcell" description:"The stemcell OS and version"`

	RequestAndWait bool          `long:"request-and-wait" description:"Request and wait for compilations to finish"`
	WaitTimeout    time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`
}

func (o *Opts) GetCompiledReleaseVersion(api *client.Client) (*models.CRVInfoResponse, error) {
	releaseRef := models.ReleaseRef{
		Name:     o.Release.Name,
		Version:  o.Release.Version,
		Checksum: models.Checksum(o.ReleaseChecksum.String()),
	}
	stemcellRef := models.StemcellRef{
		OS:      o.Stemcell.OS,
		Version: o.Stemcell.Version,
	}

	if o.RequestAndWait {
		return client.RequestAndWait(api, releaseRef, stemcellRef)
	}

	return api.CompiledReleaseVersionInfo(models.CRVInfoRequest{
		Data: models.CRVInfoRequestData{
			Release:  releaseRef,
			Stemcell: stemcellRef,
		},
	})
}
