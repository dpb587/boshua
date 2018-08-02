package graphql

import (
	artifactgraphql "github.com/dpb587/boshua/artifact/graphql"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/graphql-go/graphql"
)

var osField = &graphql.Field{
	Type: graphql.String,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		if source, ok := p.Source.(compilation.Artifact); ok {
			return source.OS.Name, nil
		}

		return nil, nil
	},
}

var versionField = &graphql.Field{
	Type: graphql.String,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		if source, ok := p.Source.(compilation.Artifact); ok {
			return source.OS.Version, nil
		}

		return nil, nil
	},
}

var labelsField = &graphql.Field{
	Type: graphql.NewList(graphql.String),
}

var tarballField = &graphql.Field{
	Type: artifactgraphql.ArtifactType,
}
