package graphql

import (
	artifactgraphql "github.com/dpb587/boshua/artifact/graphql"
	"github.com/graphql-go/graphql"
)

var osField = &graphql.Field{
	Type: graphql.String,
}

var versionField = &graphql.Field{
	Type: graphql.String,
}

var iaasField = &graphql.Field{
	Type: graphql.String,
}

var hypervisorField = &graphql.Field{
	Type: graphql.String,
}

var diskFormatField = &graphql.Field{
	Type: graphql.String,
}

var tarballField = &graphql.Field{
	Type: artifactgraphql.ArtifactType,
}