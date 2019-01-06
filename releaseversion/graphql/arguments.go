package graphql

import (
	"github.com/graphql-go/graphql"
)

var nameArgument = &graphql.ArgumentConfig{
	Type: graphql.String,
}

var versionArgument = &graphql.ArgumentConfig{
	Type: graphql.String,
}

var checksumArgument = &graphql.ArgumentConfig{
	Type: graphql.String,
}

var uriArgument = &graphql.ArgumentConfig{
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
