package analysis

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon/opts"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	releaseopts "github.com/dpb587/boshua/releaseversion/cli/opts"
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
	ReleaseOpts  *releaseopts.Opts
	AnalysisOpts *opts.Opts
}

func (o *CmdOpts) getAnalysis() (analysis.Artifact, error) {
	subject, err := o.ReleaseOpts.Artifact()
	if err != nil {
		return analysis.Artifact{}, errors.Wrap(err, "loading release")
	}

	analysisRef := analysis.Reference{
		Subject:  subject,
		Analyzer: o.AnalysisOpts.Analyzer,
	}

	analysisIndex, err := o.AppOpts.GetAnalysisIndex(analysisRef)
	if err != nil {
		return analysis.Artifact{}, errors.Wrap(err, "loading analysis index")
	}

	results, err := analysisIndex.Filter(analysisRef)
	if err != nil {
		return analysis.Artifact{}, errors.Wrap(err, "finding analysis")
	}

	if len(results) == 0 {
		if o.AnalysisOpts.NoWait {
			return analysis.Artifact{}, errors.New("no analysis found")
		}

		scheduler, err := o.AppOpts.GetScheduler()
		if err != nil {
			return analysis.Artifact{}, errors.Wrap(err, "loading scheduler")
		}

		err = analysisdatastore.CreateAnalysis(
			scheduler,
			analysisRef,
			[]string{
				"release",
				fmt.Sprintf("--release-name=%s", subject.Name),
				fmt.Sprintf("--release-version=%s", subject.Version),
				// TODO more options; generate from subject
			},
		)
		if err != nil {
			return analysis.Artifact{}, errors.Wrap(err, "creating analysis")
		}

		results, err = analysisIndex.Filter(analysisRef)
		if err != nil {
			return analysis.Artifact{}, errors.Wrap(err, "finding finished analysis")
		}
	}

	result, err := analysisdatastore.RequireSingleResult(results)
	if err != nil {
		return analysis.Artifact{}, errors.Wrap(err, "finding analysis")
	}

	return result, nil
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
