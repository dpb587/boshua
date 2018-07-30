package graphql

import (
	artifactgraphql "github.com/dpb587/boshua/artifact/graphql"
	"github.com/graphql-go/graphql"
)

var osField = &graphql.Field{
	Name: "StemcellOS",
	Type: graphql.String,
}

var versionField = &graphql.Field{
	Name: "StemcellVersion",
	Type: graphql.String,
}

var labelsField = &graphql.Field{
	Name: "StemcellLabels",
	Type: graphql.NewList(graphql.String),
}

var iaasField = &graphql.Field{
	Name: "StemcellIaaS",
	Type: graphql.String,
}

var hypervisorField = &graphql.Field{
	Name: "StemcellHypervisor",
	Type: graphql.String,
}

var diskFormatField = &graphql.Field{
	Name: "StemcellDiskFormat",
	Type: graphql.String,
}

var flavorField = &graphql.Field{
	Name: "StemcellFlavor",
	Type: graphql.String,
}

var tarballField = &graphql.Field{
	Name: "StemcellTarball",
	Type: artifactgraphql.ArtifactType,
}
