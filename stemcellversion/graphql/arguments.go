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

var diskFormatArgument = &graphql.ArgumentConfig{
	Type: graphql.String,
}

var flavorArgument = &graphql.ArgumentConfig{
	Type: graphql.String,
}

var labelsArgument = &graphql.ArgumentConfig{
	Type: graphql.NewList(graphql.String),
}

var limitMinArgument = &graphql.ArgumentConfig{
	Type: graphql.Float,
}

var limitMaxArgument = &graphql.ArgumentConfig{
	Type: graphql.Float,
}

var limitFirstArgument = &graphql.ArgumentConfig{
	Type: graphql.Float,
}

var limitOffsetArgument = &graphql.ArgumentConfig{
	Type: graphql.Float,
}
