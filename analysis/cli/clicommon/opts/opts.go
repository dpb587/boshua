package opts

import (
	"time"

	"github.com/dpb587/boshua/analysis"
)

type Opts struct {
	Analyzer analysis.AnalyzerName `long:"analyzer" description:"The analyzer name"`

	NoWait      bool          `long:"no-wait" description:"Do not request and wait for analysis if not already available"`
	WaitTimeout time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`
}
