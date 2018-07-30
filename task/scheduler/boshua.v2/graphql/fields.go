package graphql

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversiongraphql "github.com/dpb587/boshua/stemcellversion/graphql"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func NewStemcellAnalysisField(s scheduler.Scheduler, index datastore.Index) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewObject(
			graphql.ObjectConfig{
				Name:        "ScheduledTask",
				Description: "A scheduled task status.",
				Fields: graphql.Fields{
					"status": &graphql.Field{
						Type: graphql.String,
					},
				},
			},
		),
		Description: "Schedule stemcell analysis",
		Args:        stemcellversiongraphql.FilterArgs,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			f, err := datastore.FilterParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing args")
			}

			results, err := index.GetArtifacts(f)
			if err != nil {
				return nil, errors.Wrap(err, "finding stemcell")
			}

			result, err := datastore.RequireSingleResult(results)
			if err != nil {
				return nil, errors.Wrap(err, "finding stemcell")
			}

			analysisRef := analysis.Reference{
				Subject:  result,
				Analyzer: analysis.AnalyzerName("stemcellpackages.v1"), // TODO arg
			}

			task, err := s.ScheduleAnalysis(analysisRef)
			if err != nil {
				return nil, errors.Wrap(err, "scheduling task")
			}

			status, err := task.Status()
			if err != nil {
				return nil, errors.Wrap(err, "checking status")
			}

			return map[string]interface{}{
				"status": status,
			}, nil
		},
	}
}
