package opts

import (
	"time"

	"github.com/dpb587/boshua/api/v2/client"
	"github.com/dpb587/boshua/api/v2/models"
	"github.com/dpb587/boshua/cli/client/args"
)

type Opts struct {
	Release         args.Release  `long:"release" description:"The release name and version"`
	ReleaseChecksum args.Checksum `long:"release-checksum" description:"The release checksum"`
	OS              args.OS       `long:"os" description:"The OS and version"`

	RequestAndWait bool          `long:"request-and-wait" description:"Request and wait for compilations to finish"`
	WaitTimeout    time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`
}

func (o *Opts) GetCompiledReleaseVersion(api *client.Client) (*models.CRVInfoResponse, error) {
	releaseVersionRef := models.ReleaseVersionRef{
		Name:     o.Release.Name,
		Version:  o.Release.Version,
		Checksum: o.ReleaseChecksum.ImmutableChecksum,
	}
	osVersionRef := models.OSVersionRef{
		Name:    o.OS.Name,
		Version: o.OS.Version,
	}

	if o.RequestAndWait {
		return client.RequestAndWait(api, releaseVersionRef, osVersionRef)
	}

	return api.CompiledReleaseVersionInfo(models.CRVInfoRequest{
		Data: models.CRVInfoRequestData{
			ReleaseVersionRef: releaseVersionRef,
			OSVersionRef:      osVersionRef,
		},
	})
}
