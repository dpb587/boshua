package graphql

import (
	analysisgraphql "github.com/dpb587/boshua/analysis/graphql"
	artifactgraphql "github.com/dpb587/boshua/artifact/graphql"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func newStemcellAnalysis(index datastore.AnalysisIndex) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "StemcellAnalysis",
			Description: "Analysis results of a stemcell.",
			Fields: graphql.Fields{
				"results": analysisgraphql.NewResultsField("Stemcell", index),
				// "stemcellmanifestV1": github.com/dpb587/boshua/stemcellversion/analyzers/stemcellmanifest.v1/graphql.NewField(index),
			},
		},
	)
}

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

						result, err := datastore.GetArtifact(index, f)
						if err != nil {
							return nil, errors.Wrap(err, "finding light stemcell")
						}

						return result.MetalinkFile(), nil
					},
				},
				"tarball": tarballField,
				"analyzers": &graphql.Field{
					Name: "StemcellAnalyzersField",
					Type: graphql.NewList(graphql.String),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if source, ok := p.Source.(stemcellversion.Artifact); ok {
							return source.SupportedAnalyzers(), nil
						}

						return nil, nil
					},
				},
				"analysis": &graphql.Field{
					Name: "StemcellAnalysisField",
					Type: newStemcellAnalysis(index.(datastore.AnalysisIndex)), // TODO unsafe
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						// better way?
						return p.Source, nil
					},
				},
			},
		},
	)
}

var ListedStemcell = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "StemcellSummary",
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
