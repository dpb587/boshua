package analysis

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon/opts"
	"github.com/dpb587/boshua/analysis/cli/cliutil"
	"github.com/dpb587/boshua/config/provider/setter"
	cmdopts "github.com/dpb587/boshua/main/boshua/cmd/opts"
	pivnetfileopts "github.com/dpb587/boshua/pivnetfile/cli/opts"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
)

type Cmd struct {
	setter.AppConfig `no-flag:"true"`
	*opts.Opts

	ArtifactCmd     ArtifactCmd     `command:"artifact" description:"For showing the analysis artifact"`
	DownloadCmd     DownloadCmd     `command:"download" description:"For downloading the analysis locally"`
	ResultsCmd      ResultsCmd      `command:"results" description:"For showing the results of an analysis"`
	StoreResultsCmd StoreResultsCmd `command:"store-results" description:"For storing the results of an analysis"`
}

func (c *Cmd) Execute(extra []string) error {
	c.ArtifactCmd.AppConfig = c.AppConfig
	return c.ArtifactCmd.Execute(extra)
}

type CmdOpts struct {
	AppOpts        *cmdopts.Opts `no-flag:"true"`
	PivnetFileOpts *pivnetfileopts.Opts
	AnalysisOpts   *opts.Opts
}

func (o *CmdOpts) getAnalysis() (analysis.Artifact, error) {
	cfgProvider, err := o.AppOpts.GetConfig()
	if err != nil {
		return analysis.Artifact{}, errors.Wrap(err, "loading config")
	}

	return cliutil.LoadAnalysis(
		cfgProvider.GetPivnetFileAnalysisIndex,
		func() (analysis.Subject, error) {
			return o.PivnetFileOpts.Artifact(cfgProvider)
		},
		o.AnalysisOpts,
		cfgProvider.GetScheduler,
		schedulerpkg.DefaultStatusChangeCallback,
	)
}

func New(app *cmdopts.Opts, pivnetfile *pivnetfileopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmdOpts := &CmdOpts{
		AppOpts:        app,
		PivnetFileOpts: pivnetfile,
		AnalysisOpts:   cmd.Opts,
	}

	cmd.ArtifactCmd.CmdOpts = cmdOpts
	cmd.DownloadCmd.CmdOpts = cmdOpts
	cmd.ResultsCmd.CmdOpts = cmdOpts
	cmd.StoreResultsCmd.CmdOpts = cmdOpts

	return cmd
}
