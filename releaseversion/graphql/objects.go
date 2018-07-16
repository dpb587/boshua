package graphql

import (
	"github.com/dpb587/boshua/releaseversion"
	"github.com/graphql-go/graphql"
)

var Release = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Release",
		Description: "A specific version of a release.",
		Fields: graphql.Fields{
			"name":           nameField,
			"version":        versionField,
			"source_tarball": sourceTarballField,
			"analyzers": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if source, ok := p.Source.(releaseversion.Artifact); ok {
						return source.SupportedAnalyzers(), nil
					}

					return nil, nil
				},
			},
		},
	},
)

var ListedRelease = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Release",
		Description: "A specific version of a release.",
		Fields: graphql.Fields{
			"name":           nameField,
			"version":        versionField,
			"source_tarball": sourceTarballField,
		},
	},
)
