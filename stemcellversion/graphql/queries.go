package graphql

import (
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func NewFilterArgs() graphql.FieldConfigArgument {
	return graphql.FieldConfigArgument{
		"os":          osArgument,
		"version":     versionArgument,
		"iaas":        iaasArgument,
		"hypervisor":  hypervisorArgument,
		"diskFormat":  diskFormatArgument,
		"flavor":      flavorArgument,
		"labels":      labelsArgument,
		"limitMin":    limitMinArgument,
		"limitMax":    limitMaxArgument,
		"limitFirst":  limitFirstArgument,
		"limitOffset": limitOffsetArgument,
	}
}

func NewListQuery(index datastore.Index, analysisGetter analysisdatastore.NamedGetter) *graphql.Field {
	return &graphql.Field{
		Name: "StemcellListQuery",
		Type: graphql.NewList(newStemcellObject(index, analysisGetter)),
		Args: NewFilterArgs(),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			f, err := datastore.FilterParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing filter args")
			}

			l, err := datastore.LimitParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing limit args")
			}

			results, err := index.GetArtifacts(f, l)
			if err != nil {
				return nil, errors.Wrap(err, "finding stemcells")
			}

			return results, nil
		},
	}
}
