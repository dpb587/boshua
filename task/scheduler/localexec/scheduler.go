package localexec

import (
	"os/exec"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/analyzer/factory"
	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	compilationtask "github.com/dpb587/boshua/releaseversion/compilation/task"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/task"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/dpb587/boshua/task/scheduler/storecommon"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	config               Config
	logger               logrus.FieldLogger
	releaseVersionIndex  releaseversiondatastore.Index
	stemcellVersionIndex stemcellversiondatastore.Index
}

var _ scheduler.Scheduler = &Scheduler{}

func New(config Config, releaseVersionIndex releaseversiondatastore.Index, stemcellVersionIndex stemcellversiondatastore.Index, logger logrus.FieldLogger) scheduler.Scheduler {
	return Scheduler{
		config:               config,
		releaseVersionIndex:  releaseVersionIndex,
		stemcellVersionIndex: stemcellVersionIndex,
		logger:               logger,
	}
}

func (s Scheduler) cmd(args ...string) *exec.Cmd {
	cmd := exec.Command(s.config.Exec, append(s.config.Args, args...)...)

	return cmd
}

func (s Scheduler) schedule(tt *task.Task, subject interface{}) (scheduler.ScheduledTask, error) {
	// definitely not parallel-safe
	return newScheduledTask(s.cmd, tt, subject), nil
}

func (s Scheduler) ScheduleAnalysis(analysisRef analysis.Reference) (scheduler.ScheduledTask, error) {
	tt, err := factory.SoonToBeDeprecatedFactory.BuildTask(analysisRef.Analyzer, analysisRef.Subject)
	if err != nil {
		return nil, errors.Wrap(err, "preparing task")
	}

	tt = storecommon.AppendAnalysisStore(tt, analysisRef)

	return s.schedule(tt, analysisRef)
}

func (s Scheduler) ScheduleCompilation(f compilationdatastore.FilterParams) (scheduler.ScheduledTask, error) {
	release, err := releaseversiondatastore.GetArtifact(s.releaseVersionIndex, f.Release)
	if err != nil {
		return nil, errors.Wrap(err, "finding release")
	}

	// TODO switch to receiving stemcell; override iaas to scheduler config settings
	stemcell, err := stemcellversiondatastore.GetArtifact(s.stemcellVersionIndex, stemcellversiondatastore.FilterParams{
		OSExpected:      true,
		OS:              f.OS.Name,
		VersionExpected: true,
		Version:         f.OS.Version,
		IaaSExpected:    true,
		IaaS:            "aws",
		FlavorExpected:  true,
		Flavor:          "light",
	})
	if err != nil {
		return nil, errors.Wrap(err, "filtering stemcell")
	}

	tt, err := compilationtask.New(release, stemcell)
	if err != nil {
		return nil, errors.Wrap(err, "preparing task")
	}

	tt = storecommon.AppendCompilationStore(tt, release, stemcell)

	return s.schedule(tt, f)
}
