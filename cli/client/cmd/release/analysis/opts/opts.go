package opts

import (
	"time"
)

type Opts struct {
	Analyzer string `long:"analyzer" description:"The analyzer name"`

	RequestAndWait bool          `long:"request-and-wait" description:"Request and wait for compilations to finish"`
	WaitTimeout    time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`
}
