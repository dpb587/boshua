package analysis

import (
	"fmt"
	"os"
	"time"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon/opts"
	"github.com/dpb587/boshua/analysis/cli/cliutil"
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	"github.com/dpb587/boshua/releaseversion"
	releaseopts "github.com/dpb587/boshua/releaseversion/cli/opts"
	"github.com/dpb587/boshua/task"
)

type Cmd struct {
	*opts.Opts

	ArtifactCmd     ArtifactCmd     `command:"artifact" description:"For showing the analysis artifact"`
	ResultsCmd      ResultsCmd      `command:"results" description:"For showing the results of an analysis"`
	StoreResultsCmd StoreResultsCmd `command:"store-results" description:"For storing the results of an analysis"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.ArtifactCmd.Execute(extra)
}

type CmdOpts struct {
	AppOpts      *cmdopts.Opts `no-flag:"true"`
	ReleaseOpts  *releaseopts.Opts
	AnalysisOpts *opts.Opts
}

func (o *CmdOpts) getAnalysis() (analysis.Artifact, error) {
	var artifact releaseversion.Artifact

	return cliutil.LoadAnalysis(
		o.AppOpts.GetAnalysisIndex,
		func() (analysis.Subject, error) {
			var err error
			artifact, err = o.ReleaseOpts.Artifact()
			return artifact, err
		},
		o.AnalysisOpts,
		o.AppOpts.GetScheduler,
		func(status task.Status) {
			// TODO normalize opts
			fmt.Fprintf(os.Stderr, "%s [%s/%s] analysis is %s\n", time.Now().Format("15:04:05"), artifact.Name, artifact.Version, status)
		},
	)
}

func New(app *cmdopts.Opts, release *releaseopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmdOpts := &CmdOpts{
		AppOpts:      app,
		ReleaseOpts:  release,
		AnalysisOpts: cmd.Opts,
	}

	cmd.ArtifactCmd.CmdOpts = cmdOpts
	cmd.ResultsCmd.CmdOpts = cmdOpts
	cmd.StoreResultsCmd.CmdOpts = cmdOpts

	return cmd
}
