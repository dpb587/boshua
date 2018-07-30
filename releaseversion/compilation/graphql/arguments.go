package graphql

import (
	"github.com/graphql-go/graphql"
)

// TODO technically should support IaaS-specifics?

var osArgument = &graphql.ArgumentConfig{
	Type: graphql.String,
}

var versionArgument = &graphql.ArgumentConfig{
	Type: graphql.String,
}

var labelsArgument = &graphql.ArgumentConfig{
	Type: graphql.NewList(graphql.String),
}
