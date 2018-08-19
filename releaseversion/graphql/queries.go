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
		"name":     nameArgument,
		"version":  versionArgument,
		"checksum": checksumArgument,
		"uri":      uriArgument,
		"labels":   labelsArgument,
	}
}

func NewListQuery(index datastore.Index) *graphql.Field {
	return &graphql.Field{
		Name: "ReleaseListQuery",
		Type: graphql.NewList(ListedRelease),
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

// TODO compilation should be optional
func NewQuery(index datastore.Index, compilationIndex compilationdatastore.Index, analysisGetter analysisdatastore.NamedGetter) *graphql.Field {
	return &graphql.Field{
		Name: "ReleaseQuery",
		Type: newReleaseObject(index, compilationIndex, analysisGetter),
		Args: NewFilterArgs(),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			f, err := datastore.FilterParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing args")
			}

			result, err := datastore.GetArtifact(index, f)
			if err != nil {
				return nil, errors.Wrap(err, "finding release")
			}

			return result, err
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
