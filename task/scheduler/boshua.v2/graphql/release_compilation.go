package graphql

import (
	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	releaseversiongraphql "github.com/dpb587/boshua/releaseversion/graphql"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func NewReleaseCompilationField(s scheduler.Scheduler, releaseVersionIndex releaseversiondatastore.Index, stemcellVersionIndex stemcellversiondatastore.Index, releaseCompilationIndex compilationdatastore.Index) *graphql.Field {
	args := releaseversiongraphql.NewFilterArgs()
	// TODO support stemcell precision; objects
	args["osName"] = &graphql.ArgumentConfig{
		Type: graphql.String,
	}
	args["osVersion"] = &graphql.ArgumentConfig{
		Type: graphql.String,
	}

	return &graphql.Field{
		Type:        scheduledTask,
		Description: "Schedule release analysis",
		Args:        args,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			releaseFilter, err := releaseversiondatastore.FilterParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing args")
			}

			scheduledTask, err := s.ScheduleCompilation(compilationdatastore.FilterParams{
				Release: releaseFilter,
				OS: osversiondatastore.FilterParams{
					NameExpected:    true,
					Name:            p.Args["osName"].(string), // TODO err checking
					VersionExpected: true,
					Version:         p.Args["osVersion"].(string), // TODO err checking
				},
			})
			if err != nil {
				return nil, errors.Wrap(err, "scheduling task")
			}

			status, err := scheduledTask.Status()
			if err != nil {
				return nil, errors.Wrap(err, "checking status")
			}

			if status == scheduler.StatusSucceeded {
				// TODO better way to avoid repeated flushes?
				err = releaseCompilationIndex.FlushCompilationCache()
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

func NewReleaseCompilationAnalysisField(s scheduler.Scheduler, index compilationdatastore.Index) *graphql.Field {
	args := releaseversiongraphql.NewFilterArgs()
	// TODO support stemcell precision; objects
	args["osName"] = &graphql.ArgumentConfig{
		Type: graphql.String,
	}
	args["osVersion"] = &graphql.ArgumentConfig{
		Type: graphql.String,
	}
	args["analyzer"] = &graphql.ArgumentConfig{
		Type: graphql.String,
	}

	return &graphql.Field{
		Type:        scheduledTask,
		Description: "Schedule release compilation analysis",
		Args:        args,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			analyzer, ok := p.Args["analyzer"].(string)
			if !ok {
				return nil, errors.New("parsing args: analyzer: invalid")
			}

			f, err := releaseversiondatastore.FilterParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing args")
			}

			result, err := compilationdatastore.GetCompilationArtifact(index, compilationdatastore.FilterParams{
				Release: f,
				OS: osversiondatastore.FilterParams{
					NameExpected:    true,
					Name:            p.Args["osName"].(string), // TODO err checking
					VersionExpected: true,
					Version:         p.Args["osVersion"].(string), // TODO err checking
				},
			})
			if err != nil {
				return nil, errors.Wrap(err, "finding compilation")
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
				// TODO better way to avoid repeated flushes?
				if analysisIndex, ok := index.(analysisdatastore.Index); ok {
					err = analysisIndex.FlushAnalysisCache()
					if err != nil {
						return nil, errors.Wrap(err, "flushing cache")
					}
				}
			}

			return map[string]interface{}{
				"status": status,
			}, nil
		},
	}
}
