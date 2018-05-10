package analysis

import (
	"errors"
	"fmt"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/cli/clicommon/opts"
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	compiledreleaseopts "github.com/dpb587/boshua/compiledreleaseversion/cli/opts"
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
	AppOpts             *cmdopts.Opts `no-flag:"true"`
	CompiledReleaseOpts *compiledreleaseopts.Opts
	AnalysisOpts        *opts.Opts
}

func (o *CmdOpts) getAnalysis() (analysis.Artifact, error) {
	datastore, err := o.AppOpts.GetCompiledReleaseIndex("default")
	if err != nil {
		return analysis.Artifact{}, fmt.Errorf("loading compiled release index: %v", err)
	}

	_, err = datastore.Find(o.CompiledReleaseOpts.Reference())
	if err != nil {
		return analysis.Artifact{}, fmt.Errorf("finding compiled release: %v", err)
	}

	return analysis.Artifact{}, errors.New("TODO")

	// return client.RequireCompiledReleaseVersionAnalysis(
	// 	ref.ReleaseVersion,
	// 	ref.OSVersion,
	// 	analyzer,
	// 	func(task scheduler.TaskStatus) {
	// 		if o.AppOpts.Quiet {
	// 			return
	// 		}
	//
	// 		fmt.Fprintf(
	// 			os.Stderr,
	// 			"boshua | %s | requesting compiled release analysis: %s/%s: %s/%s: %s: task is %s\n",
	// 			time.Now().Format("15:04:05"),
	// 			ref.OSVersion.Name,
	// 			ref.OSVersion.Version,
	// 			ref.ReleaseVersion.Name,
	// 			ref.ReleaseVersion.Version,
	// 			analyzer,
	// 			task.Status,
	// 		)
	// 	},
	// )
}

func New(app *cmdopts.Opts, compiledrelease *compiledreleaseopts.Opts) *Cmd {
	cmd := &Cmd{
		Opts: &opts.Opts{},
	}

	cmdOpts := &CmdOpts{
		AppOpts:             app,
		CompiledReleaseOpts: compiledrelease,
		AnalysisOpts:        cmd.Opts,
	}

	cmd.ArtifactCmd.CmdOpts = cmdOpts
	cmd.ResultsCmd.CmdOpts = cmdOpts

	return cmd
}
