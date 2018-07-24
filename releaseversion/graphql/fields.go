package graphql

import (
	artifactgraphql "github.com/dpb587/boshua/artifact/graphql"
	"github.com/graphql-go/graphql"
)

var nameField = &graphql.Field{
	Type: graphql.String,
}

var versionField = &graphql.Field{
	Type: graphql.String,
}

var labelsField = &graphql.Field{
	Type: graphql.NewList(graphql.String),
}

var sourceTarballField = &graphql.Field{
	Type: artifactgraphql.ArtifactType,
}
