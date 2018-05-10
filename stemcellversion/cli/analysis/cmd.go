package analysis

import (
	"errors"
	"fmt"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon/opts"
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	stemcellopts "github.com/dpb587/boshua/stemcellversion/cli/opts"
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
	datastore, err := o.AppOpts.GetStemcellIndex("default")
	if err != nil {
		return analysis.Artifact{}, fmt.Errorf("loading stemcell index: %v", err)
	}

	_, err = datastore.Find(o.StemcellOpts.Reference())
	if err != nil {
		return analysis.Artifact{}, fmt.Errorf("finding stemcell: %v", err)
	}

	return analysis.Artifact{}, errors.New("TODO")

	// return client.RequireStemcellVersionAnalysis(
	// 	ref,
	// 	analyzer,
	// 	func(task scheduler.TaskStatus) {
	// 		if o.AppOpts.Quiet {
	// 			return
	// 		}
	//
	// 		fmt.Fprintf(
	// 			os.Stderr,
	// 			"boshua | %s | requesting stemcell analysis: %s/%s: %s: task is %s\n",
	// 			time.Now().Format("15:04:05"),
	// 			ref.Name(),
	// 			ref.Version,
	// 			analyzer,
	// 			task.Status,
	// 		)
	// 	},
	// )
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
