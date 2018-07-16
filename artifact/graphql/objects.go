package graphql

import (
	"github.com/dpb587/metalink"
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
				Args: graphql.FieldConfigArgument{
					"types": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(metalink.File)
					if !ok {
						return source, nil
					}

					hashTypes, hashTypesExpected := p.Args["types"].([]interface{})
					if !hashTypesExpected {
						return source.Hashes, nil
					}

					filtered := []metalink.Hash{}

					for _, hashType := range hashTypes {
						for _, hash := range source.Hashes {
							if hashType == hash.Type {
								filtered = append(filtered, hash)

								break
							}
						}
					}

					return filtered, nil
				},
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
