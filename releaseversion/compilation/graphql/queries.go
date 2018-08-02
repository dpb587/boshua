package graphql

import (
	"fmt"

	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

var FilterArgs = graphql.FieldConfigArgument{
	"os":      osArgument,
	"version": versionArgument,
	"labels":  labelsArgument, // TODO unsupport?
}

func NewQuery(index datastore.Index) *graphql.Field {
	return &graphql.Field{
		Name: "ReleaseCompilationQuery",
		Type: graphql.NewList(newCompilationObject(index)),
		Args: FilterArgs,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			release, ok := p.Source.(releaseversion.Artifact)
			if !ok {
				panic(fmt.Errorf("not a release: %#+v", p.Source))
			}

			releaseParams := releaseversiondatastore.FilterParamsFromArtifact(release)

			osParams, err := osversiondatastore.FilterParamsFromMap(p.Args)
			if err != nil {
				return nil, errors.Wrap(err, "parsing args")
			}

			results, err := index.GetCompilationArtifacts(datastore.FilterParams{
				Release: releaseParams,
				OS:      osParams,
			})
			if err != nil {
				return nil, errors.Wrap(err, "finding compilation")
			}

			return results, err
		},
	}
}
