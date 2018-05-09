package opts

import (
	"github.com/dpb587/boshua/cli/args"
	"github.com/dpb587/boshua/releaseversion"
)

type Opts struct {
	Release         args.Release   `long:"release" description:"The release name and version"`
	ReleaseChecksum *args.Checksum `long:"release-checksum" description:"The release checksum"`
}

func (o Opts) Reference() releaseversion.Reference {
	ref := releaseversion.Reference{
		Name:    o.Release.Name,
		Version: o.Release.Version,
	}

	if o.ReleaseChecksum != nil {
		ref.Checksums = append(ref.Checksums, o.ReleaseChecksum.ImmutableChecksum)
	}

	return ref
}
