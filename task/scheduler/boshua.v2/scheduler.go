package boshuaV2

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/analysis"
	boshuaV2 "github.com/dpb587/boshua/artifact/datastore/datastoreutil/boshua.v2"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	releaseversiongraphql "github.com/dpb587/boshua/releaseversion/graphql"
	"github.com/dpb587/boshua/stemcellversion"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversiongraphql "github.com/dpb587/boshua/stemcellversion/graphql"
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
	switch subject := analysisRef.Subject.(type) {
	case stemcellversion.Artifact:
		return s.scheduleStemcellAnalysis(subject, analysisRef.Analyzer)
	case releaseversion.Artifact:
		return s.scheduleReleaseAnalysis(subject, analysisRef.Analyzer)
	case compilation.Artifact:
		panic("TODO")
	default:
		panic(errors.New("unsupported analysis subject")) // TODO panic?
	}
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

			var resp mutationScheduleAnalysis

			err := s.client.Execute(req, &resp)
			if err != nil {
				return task.StatusUnknown, errors.Wrap(err, "executing remote request")
			}

			return resp.Status(), nil
		},
	), nil
}

func (s scheduler) scheduleStemcellAnalysis(subject stemcellversion.Artifact, analyzer analysis.AnalyzerName) (task.ScheduledTask, error) {
	mutationFilter, mutationVarsTypes, mutationVars := stemcellversiongraphql.BuildListQueryArgs(stemcellversiondatastore.FilterParamsFromArtifact(subject))
	mutationVars["analyzer"] = analyzer

	return s.schedule(
		fmt.Sprintf(`mutation _(%s, $analyzer: String!) { scheduleStemcellAnalysis(%s, analyzer: $analyzer) { status } }`, mutationVarsTypes, mutationFilter),
		mutationVars,
	)
}

func (s scheduler) scheduleReleaseAnalysis(subject releaseversion.Artifact, analyzer analysis.AnalyzerName) (task.ScheduledTask, error) {
	mutationFilter, mutationVarsTypes, mutationVars := releaseversiongraphql.BuildListQueryArgs(releaseversiondatastore.FilterParamsFromArtifact(subject))
	mutationVars["analyzer"] = analyzer

	return s.schedule(
		fmt.Sprintf(`mutation _(%s, $analyzer: String!) { scheduleReleaseAnalysis(%s, analyzer: $analyzer) { status } }`, mutationVarsTypes, mutationFilter),
		mutationVars,
	)
}
