package graphql

import (
	"github.com/graphql-go/graphql"
)

var scheduledTask = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "ScheduledTask",
		Description: "A scheduled task status.",
		Fields: graphql.Fields{
			"status": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
