package analysis

import (
	"fmt"
	"os"
	"time"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon/opts"
	"github.com/dpb587/boshua/analysis/cli/cliutil"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	"github.com/dpb587/boshua/stemcellversion"
	stemcellopts "github.com/dpb587/boshua/stemcellversion/cli/opts"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
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
	StemcellOpts *stemcellopts.Opts
	AnalysisOpts *opts.Opts
}

func (o *CmdOpts) getAnalysis() (analysis.Artifact, error) {
	var artifact stemcellversion.Artifact

	return cliutil.LoadAnalysis(
		func(analysis.Reference) (analysisdatastore.Index, error) {
			idx, err := o.StemcellOpts.Index("default")
			if err != nil {
				return nil, errors.Wrap(err, "getting index")
			}

			analysisIdx, analysisSupported := idx.(analysisdatastore.Index)
			if !analysisSupported {
				return nil, errors.Wrap(err, "getting analysis index")
			}

			return analysisIdx, nil
		},
		func() (analysis.Subject, error) {
			var err error
			artifact, err = o.StemcellOpts.Artifact()
			return artifact, err
		},
		o.AnalysisOpts,
		o.AppOpts.GetScheduler,
		func(status schedulerpkg.Status) {
			if o.AppOpts.Quiet {
				return
			}

			fmt.Fprintf(os.Stderr, "%s [%s/%s] analysis is %s\n", time.Now().Format("15:04:05"), artifact.FullName(), artifact.Version, status)
		},
	)
}

func New(app *cmdopts.Opts, stemcell *stemcellopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmdOpts := &CmdOpts{
		AppOpts:      app,
		StemcellOpts: stemcell,
		AnalysisOpts: cmd.Opts,
	}

	cmd.ArtifactCmd.CmdOpts = cmdOpts
	cmd.ResultsCmd.CmdOpts = cmdOpts
	cmd.StoreResultsCmd.CmdOpts = cmdOpts

	return cmd
}
