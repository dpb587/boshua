package opts

import (
	"time"

	"github.com/dpb587/boshua/cli/args"
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/osversion"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	"github.com/dpb587/boshua/releaseversion"
	releaseversionopts "github.com/dpb587/boshua/releaseversion/cli/opts"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
)

type Opts struct {
	Release         releaseversionopts.Release `long:"release" description:"The release name and version"`
	ReleaseChecksum *args.Checksum             `long:"release-checksum" description:"The release checksum"`

	OS args.OS `long:"os" description:"The OS and version"`

	NoWait      bool          `long:"no-wait" description:"Do not request and wait for compilation if not already available"`
	WaitTimeout time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`
}

func (o Opts) FilterParams() *datastore.FilterParams {
	return &datastore.FilterParams{
		Release: &releaseversiondatastore.FilterParams{
			NameExpected:    o.Release.Name != "",
			Name:            o.Release.Name,
			VersionExpected: o.Release.Version != "",
			Version:         o.Release.Version,
		},
		OS: &osversiondatastore.FilterParams{
			NameExpected:    o.OS.Name != "",
			Name:            o.OS.Name,
			VersionExpected: o.OS.Version != "",
			Version:         o.OS.Version,
		},
	}
}

func (o Opts) Reference() compiledreleaseversion.Reference {
	releaseVersionRef := releaseversion.Reference{
		Name:    o.Release.Name,
		Version: o.Release.Version,
	}

	if o.ReleaseChecksum != nil {
		releaseVersionRef.Checksums = append(releaseVersionRef.Checksums, o.ReleaseChecksum.ImmutableChecksum)
	}

	osVersionRef := osversion.Reference{
		Name:    o.OS.Name,
		Version: o.OS.Version,
	}

	ref := compiledreleaseversion.Reference{
		ReleaseVersion: releaseVersionRef,
		OSVersion:      osVersionRef,
	}

	return ref
}
