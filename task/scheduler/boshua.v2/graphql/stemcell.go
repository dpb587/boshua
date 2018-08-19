package graphql

import (
	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversiongraphql "github.com/dpb587/boshua/stemcellversion/graphql"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func NewStemcellAnalysisField(s scheduler.Scheduler, index datastore.Index, analysisIndexGetter analysisdatastore.NamedGetter) *graphql.Field {
	args := stemcellversiongraphql.NewFilterArgs()
	args["analyzer"] = &graphql.ArgumentConfig{
		Type: graphql.String,
	}

	return &graphql.Field{
		Type:        scheduledTask,
		Description: "Schedule stemcell analysis",
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

			result, err := datastore.GetArtifact(index, f)
			if err != nil {
				return nil, errors.Wrap(err, "finding stemcell")
			}

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
