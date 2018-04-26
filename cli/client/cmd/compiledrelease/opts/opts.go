package opts

import (
	"time"

	"github.com/dpb587/boshua/api/v2/client"
	"github.com/dpb587/boshua/api/v2/models/compiledreleaseversion"
	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/cli/client/args"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
)

type Opts struct {
	Release         args.Release  `long:"release" description:"The release name and version"`
	ReleaseChecksum args.Checksum `long:"release-checksum" description:"The release checksum"`
	OS              args.OS       `long:"os" description:"The OS and version"`

	RequestAndWait bool          `long:"request-and-wait" description:"Request and wait for compilations to finish"`
	WaitTimeout    time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`
}

func (o *Opts) GetCompiledReleaseVersion(api *client.Client) (*compiledreleaseversion.GETCompilationResponse, error) {
	var get func(releaseversion.Reference, osversion.Reference) (*compiledreleaseversion.GETCompilationResponse, error) = api.GetCompiledReleaseVersionCompilation

	if o.RequestAndWait {
		get = api.RequireCompiledReleaseVersionCompilation
	}

	return get(
		releaseversion.Reference{
			Name:      o.Release.Name,
			Version:   o.Release.Version,
			Checksums: checksum.ImmutableChecksums{o.ReleaseChecksum.ImmutableChecksum},
		},
		osversion.Reference{
			Name:    o.OS.Name,
			Version: o.OS.Version,
		},
	)
}
