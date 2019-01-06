package graphql

import (
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func NewFilterArgs() graphql.FieldConfigArgument {
	return graphql.FieldConfigArgument{
		"name":        nameArgument,
		"version":     versionArgument,
		"checksum":    checksumArgument,
		"uri":         uriArgument,
		"labels":      labelsArgument,
		"limitMin":    limitMinArgument,
		"limitMax":    limitMaxArgument,
		"limitFirst":  limitFirstArgument,
		"limitOffset": limitOffsetArgument,
	}
}

// TODO compilation should be optional
func NewListQuery(index datastore.Index, compilationIndex compilationdatastore.Index, analysisGetter analysisdatastore.NamedGetter) *graphql.Field {
	return &graphql.Field{
		Name: "ReleaseListQuery",
		Type: graphql.NewList(newReleaseObject(index, compilationIndex, analysisGetter)),
		Args: NewFilterArgs(),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			f, err := datastore.FilterParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing args")
			}

			l, err := datastore.LimitParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing limit args")
			}

			return index.GetArtifacts(f, l)
		},
	}
}

func NewLabelsQuery(r datastore.Index) *graphql.Field {
	return &graphql.Field{
		Name: "ReleaseLabelsQuery",
		Type: graphql.NewList(ReleaseLabel),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return r.GetLabels()
		},
	}
}
