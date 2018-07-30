package graphql

import (
	analysisgraphql "github.com/dpb587/boshua/analysis/graphql"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	"github.com/graphql-go/graphql"
)

func newCompilationAnalysis(index datastore.AnalysisIndex) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "ReleaseCompilationAnalysis",
			Description: "Analysis results of a release version compilation.",
			Fields: graphql.Fields{
				"results": analysisgraphql.NewResultsField("ReleaseCompilation", index),
				// "stemcellmanifestV1": github.com/dpb587/boshua/stemcellversion/analyzers/stemcellmanifest.v1/graphql.NewField(index),
			},
		},
	)
}

func newCompilationObject(index datastore.Index) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "ReleaseCompilation",
			Description: "A specific compilation of a version of a release.",
			Fields: graphql.Fields{
				"os":      osField,
				"version": versionField,
				"labels":  labelsField,
				"tarball": tarballField,
				"analyzers": &graphql.Field{
					Name: "ReleaseCompilationAnalyzers",
					Type: graphql.NewList(graphql.String),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if source, ok := p.Source.(compilation.Artifact); ok {
							return source.SupportedAnalyzers(), nil
						}

						return nil, nil
					},
				},
				"analysis": &graphql.Field{
					Name: "ReleaseCompilationAnalysisField",
					Type: newCompilationAnalysis(index.(datastore.AnalysisIndex)), // TODO unsafe
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						// better way?
						return p.Source, nil
					},
				},
			},
		},
	)
}
