package graphql

import (
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func NewFilterArgs() graphql.FieldConfigArgument {
	return graphql.FieldConfigArgument{
		"os":         osArgument,
		"version":    versionArgument,
		"iaas":       iaasArgument,
		"hypervisor": hypervisorArgument,
		"diskFormat": diskFormatArgument,
		"flavor":     flavorArgument,
		"labels":     labelsArgument,
	}
}

func NewListQuery(index datastore.Index) *graphql.Field {
	return &graphql.Field{
		Name: "StemcellListQuery",
		Type: graphql.NewList(ListedStemcell),
		Args: NewFilterArgs(),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			f, err := datastore.FilterParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing args")
			}

			return index.GetArtifacts(f)
		},
	}
}

func NewQuery(index datastore.Index, analysisGetter analysisdatastore.NamedGetter) *graphql.Field {
	return &graphql.Field{
		Name: "StemcellQuery",
		Type: newStemcellObject(index, analysisGetter),
		Args: NewFilterArgs(),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			f, err := datastore.FilterParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing args")
			}

			result, err := datastore.GetArtifact(index, f)
			if err != nil {
				return nil, errors.Wrap(err, "finding stemcell")
			}

			return result, nil
		},
	}
}
