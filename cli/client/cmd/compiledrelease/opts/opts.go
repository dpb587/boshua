package opts

import (
	"time"

	"github.com/dpb587/boshua/cli/client/args"
)

type Opts struct {
	Release         args.Release  `long:"release" description:"The release name and version"`
	ReleaseChecksum args.Checksum `long:"release-checksum" description:"The release checksum"`
	OS              args.OS       `long:"os" description:"The OS and version"`

	NoWait      bool          `long:"no-wait" description:"Do not request and wait for compilation if not already available"`
	WaitTimeout time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`
}
