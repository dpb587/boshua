package graphql

import (
	"github.com/graphql-go/graphql"
)

var osArgument = &graphql.ArgumentConfig{
	Type: graphql.String,
}

var versionArgument = &graphql.ArgumentConfig{
	Type: graphql.String,
}

var iaasArgument = &graphql.ArgumentConfig{
	Type: graphql.String,
}

var hypervisorArgument = &graphql.ArgumentConfig{
	Type: graphql.String,
}

var lightArgument = &graphql.ArgumentConfig{
	Type: graphql.Boolean,
}

var diskFormatArgument = &graphql.ArgumentConfig{
	Type: graphql.String,
}