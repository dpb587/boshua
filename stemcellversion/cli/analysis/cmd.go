package analysis

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon/opts"
	"github.com/dpb587/boshua/analysis/cli/cliutil"
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	stemcellopts "github.com/dpb587/boshua/stemcellversion/cli/opts"
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
	return cliutil.LoadAnalysis(
		o.AppOpts.GetAnalysisIndex,
		func() (analysis.Subject, error) {
			return o.StemcellOpts.Artifact()
		},
		o.AnalysisOpts,
		o.AppOpts.GetScheduler,
		[]string{
			"stemcell",
			fmt.Sprintf("--stemcell-os=%s", o.StemcellOpts.OS),
			fmt.Sprintf("--stemcell-version=%s", o.StemcellOpts.Version),
			fmt.Sprintf("--stemcell-iaas=%s", o.StemcellOpts.IaaS),
			fmt.Sprintf("--stemcell-hypervisor=%s", o.StemcellOpts.Hypervisor),
			fmt.Sprintf("--stemcell-flavor=%s", o.StemcellOpts.Flavor),
			// TODO more options; generate from subject
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
