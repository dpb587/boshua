package graphql

import (
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/graphql-go/graphql"
)

var Stemcell = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Stemcell",
		Description: "A specific version of a stemcell.",
		Fields: graphql.Fields{
			"os":         osField,
			"version":    versionField,
			"iaas":       iaasField,
			"hypervisor": hypervisorField,
			"diskFormat": diskFormatField,
			"tarball":    tarballField,
			"analyzers": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if source, ok := p.Source.(stemcellversion.Artifact); ok {
						return source.SupportedAnalyzers(), nil
					}

					return nil, nil
				},
			},
		},
	},
)

var ListedStemcell = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Stemcell",
		Description: "A specific version of a stemcell.",
		Fields: graphql.Fields{
			"os":         osField,
			"version":    versionField,
			"iaas":       iaasField,
			"hypervisor": hypervisorField,
			"diskFormat": diskFormatField,
			"tarball":    tarballField,
		},
	},
)
