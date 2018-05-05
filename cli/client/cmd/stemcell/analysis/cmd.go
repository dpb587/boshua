package analysis

import (
	"fmt"
	"os"
	"time"

	"github.com/dpb587/boshua/api/v2/models/analysis"
	"github.com/dpb587/boshua/api/v2/models/scheduler"
	"github.com/dpb587/boshua/cli/client/cmd/analysisutil/opts"
	cmdopts "github.com/dpb587/boshua/cli/client/cmd/opts"
	stemcellopts "github.com/dpb587/boshua/cli/client/cmd/stemcell/opts"
	"github.com/dpb587/boshua/stemcellversion"
)

type Cmd struct {
	*opts.Opts

	ArtifactCmd ArtifactCmd `command:"metalink" description:"For showing a metalink of the analysis"`
	ResultsCmd  ResultsCmd  `command:"results" description:"For showing the results of an analysis"`
}

func (c *Cmd) Execute(extra []string) error {
	return c.ArtifactCmd.Execute(extra)
}

type CmdOpts struct {
	AppOpts      *cmdopts.Opts `no-flag:"true"`
	stemcellOpts *stemcellopts.Opts
	AnalysisOpts *opts.Opts
}

func (o *CmdOpts) getAnalysis() (*analysis.GETInfoResponse, error) {
	client := o.AppOpts.GetClient()

	ref := stemcellversion.Reference{
		IaaS:       o.stemcellOpts.Stemcell.IaaS,
		Hypervisor: o.stemcellOpts.Stemcell.Hypervisor,
		OS:         o.stemcellOpts.Stemcell.OS,
		Version:    o.stemcellOpts.Stemcell.Version,
	}
	analyzer := o.AnalysisOpts.Analyzer

	if o.AnalysisOpts.NoWait {
		return client.GetStemcellVersionAnalysis(ref, analyzer)
	}

	return client.RequireStemcellVersionAnalysis(
		ref,
		analyzer,
		func(task scheduler.TaskStatus) {
			if o.AppOpts.Quiet {
				return
			}

			fmt.Fprintf(os.Stderr, "boshua | %s | requesting stemcell analysis: %s/%s: %s: task is %s\n", time.Now().Format("15:04:05"), ref.Name(), ref.Version, analyzer, task.Status)
		},
	)
}

func New(app *cmdopts.Opts, stemcell *stemcellopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmdOpts := &CmdOpts{
		AppOpts:      app,
		stemcellOpts: stemcell,
		AnalysisOpts: cmd.Opts,
	}

	cmd.ArtifactCmd.CmdOpts = cmdOpts
	cmd.ResultsCmd.CmdOpts = cmdOpts

	return cmd
}
