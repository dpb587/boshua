package graphql

import (
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func NewListQuery(r datastore.Index) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(ListedStemcell),
		Args: graphql.FieldConfigArgument{
			"os":         osArgument,
			"version":    versionArgument,
			"iaas":       iaasArgument,
			"hypervisor": hypervisorArgument,
			"diskFormat": diskFormatArgument,
			"flavor":     flavorArgument,
			"labels":     labelsArgument,
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			f, err := datastore.FilterParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing args")
			}

			return r.Filter(f)
		},
	}
}
