package opts

import (
	"github.com/dpb587/boshua/cli/client/args"
)

type Opts struct {
	Stemcell args.Stemcell `long:"stemcell" description:"The stemcell name and version"`
}
