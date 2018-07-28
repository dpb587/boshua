package graphql

import (
	artifactgraphql "github.com/dpb587/boshua/artifact/graphql"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func newStemcellObject(index datastore.Index) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Stemcell",
			Description: "A specific version of a stemcell.",
			Fields: graphql.Fields{
				"os":         osField,
				"version":    versionField,
				"labels":     labelsField,
				"iaas":       iaasField,
				"hypervisor": hypervisorField,
				"diskFormat": diskFormatField,
				"flavor":     flavorField,
				"light_tarball": &graphql.Field{
					Type: artifactgraphql.ArtifactType,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						source, ok := p.Source.(stemcellversion.Artifact)
						if !ok {
							return nil, nil
						} else if source.Flavor == "light" {
							// already light
							return nil, nil
						}

						f := datastore.FilterParamsFromArtifact(source)
						f.Flavor = "light"

						results, err := index.Filter(f)
						if err != nil {
							return nil, errors.Wrap(err, "finding light stemcell")
						}

						result, err := datastore.RequireSingleResult(results)
						if err != nil {
							return nil, errors.Wrap(err, "finding light stemcell")
						}

						return result.MetalinkFile(), err
					},
				},
				"tarball": tarballField,
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
}

var ListedStemcell = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Stemcell",
		Description: "A specific version of a stemcell.",
		Fields: graphql.Fields{
			"os":         osField,
			"version":    versionField,
			"labels":     labelsField,
			"iaas":       iaasField,
			"hypervisor": hypervisorField,
			"diskFormat": diskFormatField,
			"flavor":     flavorField,
			"tarball":    tarballField,
		},
	},
)
