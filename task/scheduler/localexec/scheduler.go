package localexec

import (
	"os/exec"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/analyzer/factory"
	"github.com/dpb587/boshua/releaseversion"
	compilationtask "github.com/dpb587/boshua/releaseversion/compilation/task"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/task"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/dpb587/boshua/task/scheduler/storecommon"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	config Config
	logger logrus.FieldLogger
}

var _ scheduler.Scheduler = &Scheduler{}

func New(config Config, logger logrus.FieldLogger) scheduler.Scheduler {
	return Scheduler{
		config: config,
		logger: logger,
	}
}

func (s Scheduler) cmd(args ...string) *exec.Cmd {
	cmd := exec.Command(s.config.Exec, append(s.config.Args, args...)...)

	return cmd
}

func (s Scheduler) schedule(tt *task.Task) (task.ScheduledTask, error) {
	return NewTask(s.cmd, tt), nil
}

func (s Scheduler) ScheduleAnalysis(analysisRef analysis.Reference) (task.ScheduledTask, error) {
	tt, err := factory.SoonToBeDeprecatedFactory.BuildTask(analysisRef.Analyzer, analysisRef.Subject)
	if err != nil {
		return nil, errors.Wrap(err, "preparing task")
	}

	tt = storecommon.AppendAnalysisStore(tt, analysisRef)

	return s.schedule(tt)
}

func (s Scheduler) ScheduleCompilation(release releaseversion.Artifact, stemcell stemcellversion.Artifact) (task.ScheduledTask, error) {
	tt, err := compilationtask.New(release, stemcell)
	if err != nil {
		return nil, errors.Wrap(err, "preparing task")
	}

	tt = storecommon.AppendCompilationStore(tt, release, stemcell)

	return s.schedule(tt)
}
