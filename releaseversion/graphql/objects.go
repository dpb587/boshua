package graphql

import (
	"github.com/dpb587/boshua/releaseversion"
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

var Release = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Release",
		Description: "A specific version of a release.",
		Fields: graphql.Fields{
			"name":           nameField,
			"version":        versionField,
			"labels":         labelsField,
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
			"labels":         labelsField,
			"source_tarball": sourceTarballField,
		},
	},
)
