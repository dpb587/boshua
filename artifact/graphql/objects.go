package graphql

import (
	"github.com/graphql-go/graphql"
)

var ArtifactType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Artifact",
		Description: "A reference to an artifact which can be downloaded.",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"hashes": &graphql.Field{
				Type: graphql.NewList(ArtifactHashType),
			},
			"size": &graphql.Field{
				Type: graphql.Int,
			},
			"urls": &graphql.Field{
				Type: graphql.NewList(ArtifactURLType),
			},
			"metaurls": &graphql.Field{
				Type: graphql.NewList(ArtifactMetaURLType),
			},
		},
	},
)

var ArtifactHashType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Hash",
		Description: "A checksum hash which can be used to verify a download.",
		Fields: graphql.Fields{
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"hash": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var ArtifactURLType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "URL",
		Description: "A URL for downloading the asset.",
		Fields: graphql.Fields{
			"url": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var ArtifactMetaURLType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "MetaURL",
		Description: "A Meta URL for downloading the asset.",
		Fields: graphql.Fields{
			"url": &graphql.Field{
				Type: graphql.String,
			},
			"mediatype": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
