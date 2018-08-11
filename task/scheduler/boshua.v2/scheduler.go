package boshuaV2

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/analysis"
	boshuaV2 "github.com/dpb587/boshua/artifact/datastore/datastoreutil/boshua.v2"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation"
	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	releaseversiongraphql "github.com/dpb587/boshua/releaseversion/graphql"
	"github.com/dpb587/boshua/stemcellversion"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversiongraphql "github.com/dpb587/boshua/stemcellversion/graphql"
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

func (s scheduler) ScheduleAnalysis(analysisRef analysis.Reference) (schedulerpkg.ScheduledTask, error) {
	switch subject := analysisRef.Subject.(type) {
	case stemcellversion.Artifact:
		return s.scheduleStemcellAnalysis(analysisRef, subject)
	case releaseversion.Artifact:
		return s.scheduleReleaseAnalysis(analysisRef, subject)
	case compilation.Artifact:
		return s.scheduleReleaseCompilationAnalysis(analysisRef, subject)
	default:
		panic(errors.New("unsupported analysis subject")) // TODO panic?
	}
}

func (s scheduler) ScheduleCompilation(f compilationdatastore.FilterParams) (schedulerpkg.ScheduledTask, error) {
	mutationFilter, mutationVarsTypes, mutationVars := releaseversiongraphql.BuildListQueryArgs(f.Release)
	mutationVars["osName"] = f.OS.Name
	mutationVars["osVersion"] = f.OS.Version

	return s.schedule(
		f,
		fmt.Sprintf(`mutation _(%s, $osName: String!, $osVersion: String!) { scheduleReleaseCompilation(%s, osName: $osName, osVersion: $osVersion) { status } }`, mutationVarsTypes, mutationFilter),
		mutationVars,
	)
}

func (s scheduler) schedule(subject interface{}, mutationQuery string, mutationVars map[string]interface{}) (schedulerpkg.ScheduledTask, error) {
	// TODO assumes caller runs Status; not calling here avoids duplicate initial requests
	return newScheduledTask(
		func() (schedulerpkg.Status, error) {
			req := graphql.NewRequest(mutationQuery)

			for k, v := range mutationVars {
				req.Var(k, v)
			}

			var resp mutationSchedule

			err := s.client.Execute(req, &resp)
			if err != nil {
				return schedulerpkg.StatusUnknown, errors.Wrap(err, "executing remote request")
			}

			return resp.Status(), nil
		},
		subject,
	), nil
}

func (s scheduler) scheduleStemcellAnalysis(analysisRef analysis.Reference, subject stemcellversion.Artifact) (schedulerpkg.ScheduledTask, error) {
	mutationFilter, mutationVarsTypes, mutationVars := stemcellversiongraphql.BuildListQueryArgs(stemcellversiondatastore.FilterParamsFromArtifact(subject))
	mutationVars["analyzer"] = analysisRef.Analyzer

	return s.schedule(
		analysisRef,
		fmt.Sprintf(`mutation _(%s, $analyzer: String!) { scheduleStemcellAnalysis(%s, analyzer: $analyzer) { status } }`, mutationVarsTypes, mutationFilter),
		mutationVars,
	)
}

func (s scheduler) scheduleReleaseAnalysis(analysisRef analysis.Reference, subject releaseversion.Artifact) (schedulerpkg.ScheduledTask, error) {
	mutationFilter, mutationVarsTypes, mutationVars := releaseversiongraphql.BuildListQueryArgs(releaseversiondatastore.FilterParamsFromArtifact(subject))
	mutationVars["analyzer"] = analysisRef.Analyzer

	return s.schedule(
		analysisRef,
		fmt.Sprintf(`mutation _(%s, $analyzer: String!) { scheduleReleaseAnalysis(%s, analyzer: $analyzer) { status } }`, mutationVarsTypes, mutationFilter),
		mutationVars,
	)
}

func (s scheduler) scheduleReleaseCompilationAnalysis(analysisRef analysis.Reference, subject compilation.Artifact) (schedulerpkg.ScheduledTask, error) {
	mutationFilter, mutationVarsTypes, mutationVars := releaseversiongraphql.BuildListQueryArgs(releaseversiondatastore.FilterParamsFromReference(subject.Release))
	mutationVars["osName"] = subject.OS.Name
	mutationVars["osVersion"] = subject.OS.Version
	mutationVars["analyzer"] = analysisRef.Analyzer

	return s.schedule(
		analysisRef,
		fmt.Sprintf(`mutation _(%s, $osName: String!, $osVersion: String!, $analyzer: String!) { scheduleReleaseCompilationAnalysis(%s, osName: $osName, osVersion: $osVersion, analyzer: $analyzer) { status } }`, mutationVarsTypes, mutationFilter),
		mutationVars,
	)
}
