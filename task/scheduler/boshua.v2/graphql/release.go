package graphql

import (
	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore"
	releaseversiongraphql "github.com/dpb587/boshua/releaseversion/graphql"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func NewReleaseAnalysisField(s scheduler.Scheduler, index datastore.Index, analysisIndexGetter analysisdatastore.NamedGetter) *graphql.Field {
	args := releaseversiongraphql.NewFilterArgs()
	args["analyzer"] = &graphql.ArgumentConfig{
		Type: graphql.String,
	}

	return &graphql.Field{
		Type:        scheduledTask,
		Description: "Schedule release analysis",
		Args:        args,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			analyzer, ok := p.Args["analyzer"].(string)
			if !ok {
				return nil, errors.New("parsing args: analyzer: invalid")
			}

			f, err := datastore.FilterParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing args")
			}

			results, err := index.GetArtifacts(f, datastore.SingleArtifactLimitParams)
			if err != nil {
				return nil, errors.Wrap(err, "finding release")
			}

			result := results[0]

			analysisRef := analysis.Reference{
				Subject:  result,
				Analyzer: analysis.AnalyzerName(analyzer),
			}

			scheduledTask, err := s.ScheduleAnalysis(analysisRef)
			if err != nil {
				return nil, errors.Wrap(err, "scheduling task")
			}

			status, err := scheduledTask.Status()
			if err != nil {
				return nil, errors.Wrap(err, "checking status")
			}

			if status == scheduler.StatusSucceeded {
				analysisIndex, err := analysisIndexGetter(result.GetDatastoreName())
				if err != nil {
					return nil, errors.Wrap(err, "loading analysis datastore")
				}

				// TODO better way to avoid repeated flushes?
				err = analysisIndex.FlushAnalysisCache()
				if err != nil {
					return nil, errors.Wrap(err, "flushing cache")
				}
			}

			return map[string]interface{}{
				"status": status,
			}, nil
		},
	}
}
