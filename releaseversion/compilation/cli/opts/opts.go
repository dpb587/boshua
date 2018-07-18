package opts

import (
	"time"

	"github.com/dpb587/boshua/cli/args"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	releaseversionopts "github.com/dpb587/boshua/releaseversion/cli/opts"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
)

type Opts struct {
	Release *releaseversionopts.Opts `no-flag:"true"`

	OS args.OS `long:"os" description:"The OS and version"`

	NoWait      bool          `long:"no-wait" description:"Do not request and wait for compilation if not already available"`
	WaitTimeout time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`
}

func (o Opts) FilterParams() *datastore.FilterParams {
	return &datastore.FilterParams{
		Release: o.Release.FilterParams(),
		OS: &osversiondatastore.FilterParams{
			NameExpected:    o.OS.Name != "",
			Name:            o.OS.Name,
			VersionExpected: o.OS.Version != "",
			Version:         o.OS.Version,
		},
	}
}
