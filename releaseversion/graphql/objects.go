package graphql

import (
	analysisgraphql "github.com/dpb587/boshua/analysis/graphql"
	"github.com/dpb587/boshua/releaseversion"
	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	compilationgraphql "github.com/dpb587/boshua/releaseversion/compilation/graphql"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/graphql-go/graphql"
)

var ReleaseLabel = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "ReleaseLabel",
		Description: "A release label.",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if source, ok := p.Source.(string); ok {
						return source, nil
					}

					return nil, nil
				},
			},
		},
	},
)

func newReleaseAnalysis(index datastore.AnalysisIndex) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "ReleaseAnalysis",
			Description: "Analysis results of a release version.",
			Fields: graphql.Fields{
				"results": analysisgraphql.NewResultsField("Release", index),
				// "stemcellmanifestV1": github.com/dpb587/boshua/stemcellversion/analyzers/stemcellmanifest.v1/graphql.NewField(index),
			},
		},
	)
}

func newReleaseObject(index datastore.Index, compilationIndex compilationdatastore.Index) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Release",
			Description: "A specific version of a release.",
			Fields: graphql.Fields{
				"name":    nameField,
				"version": versionField,
				"labels":  labelsField,
				"tarball": tarballField,
				"analyzers": &graphql.Field{
					Name: "ReleaseAnalyzers",
					Type: graphql.NewList(graphql.String),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if source, ok := p.Source.(releaseversion.Artifact); ok {
							return source.SupportedAnalyzers(), nil
						}

						return nil, nil
					},
				},
				"analysis": &graphql.Field{
					Name: "ReleaseAnalysisField",
					Type: newReleaseAnalysis(index.(datastore.AnalysisIndex)), // TODO unsafe
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						// better way?
						return p.Source, nil
					},
				},
				// TODO compilations for multiple
				"compilations": compilationgraphql.NewQuery(compilationIndex),
			},
		},
	)
}

var ListedRelease = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "ReleaseSummary",
		Description: "A specific version of a release.",
		Fields: graphql.Fields{
			"name":    nameField,
			"version": versionField,
			"labels":  labelsField,
			"tarball": tarballField,
		},
	},
)
