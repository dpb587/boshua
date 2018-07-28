package graphql

import (
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func NewListQuery(r datastore.Index) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(ListedRelease),
		Args: graphql.FieldConfigArgument{
			"name":     nameArgument,
			"version":  versionArgument,
			"checksum": checksumArgument,
			"uri":      uriArgument,
			"labels":   labelsArgument,
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			f, err := datastore.FilterParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing args")
			}

			return r.GetArtifacts(f)
		},
	}
}

func NewLabelsQuery(r datastore.Index) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(ReleaseLabel),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return r.GetLabels()
		},
	}
}
