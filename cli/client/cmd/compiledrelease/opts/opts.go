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

	NoWait      bool          `long:"no-wait" description:"Do not request and wait for compilation if not already available"`
	WaitTimeout time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`
}

func (o *Opts) GetCompiledReleaseVersion(api *client.Client) (*compiledreleaseversion.GETCompilationResponse, error) {
	var get func(releaseversion.Reference, osversion.Reference) (*compiledreleaseversion.GETCompilationResponse, error) = api.RequireCompiledReleaseVersionCompilation

	if o.NoWait {
		get = api.GetCompiledReleaseVersionCompilation
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
