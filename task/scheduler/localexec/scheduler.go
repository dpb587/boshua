package localexec

import (
	"fmt"
	"os/exec"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/analyzer/factory"
	"github.com/dpb587/boshua/config"
	"github.com/dpb587/boshua/config/provider"
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
	config Config
	logger logrus.FieldLogger
	cfg    *provider.Config
}

var _ scheduler.Scheduler = &Scheduler{}

func New(config Config, cfg *provider.Config, logger logrus.FieldLogger) scheduler.Scheduler {
	return Scheduler{
		config: config,
		cfg:    cfg,
		logger: logger,
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
	releaseVersionIndex, err := s.cfg.GetReleaseIndex(config.DefaultName)
	if err != nil {
		return nil, errors.Wrap(err, "loading release index")
	}

	releases, err := releaseVersionIndex.GetArtifacts(f.Release, releaseversiondatastore.SingleArtifactLimitParams)
	if err != nil {
		return nil, errors.Wrap(err, "finding release")
	}

	release := releases[0]

	stemcellVersionIndex, err := s.cfg.GetStemcellIndex(config.DefaultName)
	if err != nil {
		return nil, errors.Wrap(err, "loading stemcell index")
	}

	// TODO switch to receiving stemcell; override iaas to scheduler config settings
	stemcells, err := stemcellVersionIndex.GetArtifacts(
		stemcellversiondatastore.FilterParams{
			OSExpected:      true,
			OS:              f.OS.Name,
			VersionExpected: true,
			Version:         f.OS.Version,
			IaaSExpected:    true,
			IaaS:            "aws",
			FlavorExpected:  true,
			Flavor:          "light",
		},
		stemcellversiondatastore.SingleArtifactLimitParams,
	)
	if err != nil {
		return nil, errors.Wrap(err, "filtering stemcell")
	}

	stemcell := stemcells[0]

	tt, err := compilationtask.New(release, stemcell)
	if err != nil {
		return nil, errors.Wrap(err, "preparing task")
	}

	tt = storecommon.AppendCompilationStore(tt, release, stemcell, fmt.Sprintf("internal/release/%s", release.GetDatastoreName()))

	return s.schedule(tt, f)
}
