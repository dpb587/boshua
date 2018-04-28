package opts

import (
	"github.com/dpb587/boshua/cli/client/args"
)

type Opts struct {
	Release         args.Release  `long:"release" description:"The release name and version"`
	ReleaseChecksum args.Checksum `long:"release-checksum" description:"The release checksum"`
}
