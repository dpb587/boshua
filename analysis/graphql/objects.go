package graphql

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	artifactgraphql "github.com/dpb587/boshua/artifact/graphql"
	"github.com/dpb587/metalink"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func NewResultsField(namespace string, indexGetter datastore.NamedGetter) *graphql.Field {
	return &graphql.Field{
		Name: fmt.Sprintf("%sAnalysisResults", namespace),
		Args: graphql.FieldConfigArgument{
			"analyzers": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
			},
		},
		Type: graphql.NewList(graphql.NewObject(
			graphql.ObjectConfig{
				Name:        fmt.Sprintf("%sAnalysisResult", namespace),
				Description: "A specific version of a stemcell.",
				Fields: graphql.Fields{
					"analyzer": analyzerField,
					"artifact": &graphql.Field{
						Type: artifactgraphql.ArtifactType,
					},
				},
			},
		)),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			subject, ok := p.Source.(analysis.Subject)
			if !ok {
				panic(fmt.Errorf("not an analysis subject: %#+v", p.Source))
			}

			analyzers, ok := p.Args["analyzers"].([]interface{})
			if !ok {
				return nil, errors.New("analyzers: expected slice")
			}

			index, err := indexGetter(subject.GetDatastoreName())
			if err != nil {
				return nil, errors.Wrap(err, "loading analysis index")
			}

			var results []analysisResult

			for _, analyzer := range analyzers {
				analyzer, ok := analyzer.(string)
				if !ok {
					return nil, errors.New("analyzers: expected string values")
				}

				found, err := index.GetAnalysisArtifacts(analysis.Reference{
					Analyzer: analysis.AnalyzerName(analyzer),
					Subject:  subject,
				})
				if err != nil {
					return nil, nil // TODO handle missing vs internal error
				} else if len(found) == 0 {
					continue
				}

				// TODO multiple results?
				results = append(results, analysisResult{
					Analyzer: analyzer,
					Artifact: found[0].MetalinkFile(),
				})
			}

			return results, nil
		},
	}
}

type analysisResult struct {
	Analyzer string        `json:"analyzer"`
	Artifact metalink.File `json:"artifact"`
}
