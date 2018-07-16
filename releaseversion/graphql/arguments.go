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
