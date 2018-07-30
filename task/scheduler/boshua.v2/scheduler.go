package boshuaV2

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/analysis"
	boshuaV2 "github.com/dpb587/boshua/artifact/datastore/datastoreutil/boshua.v2"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	datastoregraphql "github.com/dpb587/boshua/stemcellversion/graphql"
	"github.com/dpb587/boshua/task"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type scheduler struct {
	logger logrus.FieldLogger
	config Config
	client *boshuaV2.Client
}

var _ schedulerpkg.Scheduler = &scheduler{}

func New(config Config, logger logrus.FieldLogger) schedulerpkg.Scheduler {
	return &scheduler{
		logger: logger.WithField("build.package", reflect.TypeOf(scheduler{}).PkgPath()),
		config: config,
		client: boshuaV2.NewClient(http.DefaultClient, config.BoshuaConfig, logger),
	}
}

func (s scheduler) ScheduleAnalysis(analysisRef analysis.Reference) (task.ScheduledTask, error) {
	var mutationFilter, mutationVarsTypes string
	var mutationVars map[string]interface{}

	switch analysisSubject := analysisRef.Subject.(type) {
	case stemcellversion.Artifact:
		mutationFilter, mutationVarsTypes, mutationVars = datastoregraphql.BuildListQueryArgs(datastore.FilterParamsFromArtifact(analysisSubject))
	case releaseversion.Artifact:
		panic("TODO")
	case compilation.Artifact:
		panic("TODO")
	default:
		panic(errors.New("unsupported analysis subject")) // TODO panic?
	}

	return s.schedule(
		// TODO analysis
		fmt.Sprintf(`mutation ScheduleAnalysis(%s) { scheduleStemcellAnalysis(%s) { status } }`, mutationVarsTypes, mutationFilter),
		mutationVars,
	)
}

func (s scheduler) ScheduleCompilation(release releaseversion.Artifact, stemcell stemcellversion.Artifact) (task.ScheduledTask, error) {
	panic("TODO")
}

func (s scheduler) schedule(mutationQuery string, mutationVars map[string]interface{}) (task.ScheduledTask, error) {
	// TODO assumes caller runs Status; not calling here avoids duplicate initial requests
	return NewTask(
		func() (task.Status, error) {
			req := graphql.NewRequest(mutationQuery)

			for k, v := range mutationVars {
				req.Var(k, v)
			}

			var resp mutationScheduleStemcellAnalysis

			err := s.client.Execute(req, &resp)
			if err != nil {
				return task.StatusUnknown, errors.Wrap(err, "executing remote request")
			}

			return resp.ScheduledTask.Status, nil
		},
	), nil
}
