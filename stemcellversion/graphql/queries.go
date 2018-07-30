package graphql

import (
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

var FilterArgs = graphql.FieldConfigArgument{
	"os":         osArgument,
	"version":    versionArgument,
	"iaas":       iaasArgument,
	"hypervisor": hypervisorArgument,
	"diskFormat": diskFormatArgument,
	"flavor":     flavorArgument,
	"labels":     labelsArgument,
}

func NewListQuery(index datastore.Index) *graphql.Field {
	return &graphql.Field{
		Name: "StemcellListQuery",
		Type: graphql.NewList(ListedStemcell),
		Args: FilterArgs,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			f, err := datastore.FilterParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing args")
			}

			return index.GetArtifacts(f)
		},
	}
}

func NewQuery(index datastore.Index) *graphql.Field {
	return &graphql.Field{
		Name: "StemcellQuery",
		Type: newStemcellObject(index),
		Args: FilterArgs,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			f, err := datastore.FilterParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing args")
			}

			results, err := index.GetArtifacts(f)
			if err != nil {
				return nil, errors.Wrap(err, "finding stemcell")
			}

			result, err := datastore.RequireSingleResult(results)
			if err != nil {
				return nil, errors.Wrap(err, "finding stemcell")
			}

			return result, err
		},
	}
}
