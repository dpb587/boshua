package analysis

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon/opts"
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	stemcellopts "github.com/dpb587/boshua/stemcellversion/cli/opts"
	"github.com/pkg/errors"
)

type Cmd struct {
	*opts.Opts

	ArtifactCmd ArtifactCmd `command:"artifact" description:"For showing the analysis artifact"`
	ResultsCmd  ResultsCmd  `command:"results" description:"For showing the results of an analysis"`
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
	return analysis.Artifact{}, errors.New("TODO resurrect functionality")
	// index, err := o.AppOpts.GetStemcellIndex("default")
	// if err != nil {
	// 	return analysis.Artifact{}, errors.Wrap(err, "loading stemcell index")
	// }
	//
	// scheduler, err := o.AppOpts.GetScheduler()
	// if err != nil {
	// 	return analysis.Artifact{}, errors.Wrap(err, "loading scheduler")
	// }
	//
	// _, subject, err := datastore.FindOrCreateAnalysis(index, scheduler, o.StemcellOpts.Reference(), o.AnalysisOpts.Analyzer)
	// if err != nil {
	// 	return analysis.Artifact{}, err // intentional no Wrap
	// }
	//
	// return subject, nil
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

	return cmd
}
