package boshuaV2

import (
	"fmt"
	"net/http"
	"reflect"

	boshuaV2 "github.com/dpb587/boshua/artifact/datastore/datastoreutil/boshua.v2"
	osversiongraphql "github.com/dpb587/boshua/osversion/graphql"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiongraphql "github.com/dpb587/boshua/releaseversion/graphql"
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Index struct {
	logger logrus.FieldLogger
	config Config
	client *boshuaV2.Client
}

var _ datastore.Index = &Index{}

func New(config Config, logger logrus.FieldLogger) *Index {
	return &Index{
		logger: logger.WithField("build.package", reflect.TypeOf(Index{}).PkgPath()),
		config: config,
		client: boshuaV2.NewClient(http.DefaultClient, config.BoshuaConfig, logger),
	}
}

func (i *Index) GetCompilationArtifacts(f datastore.FilterParams) ([]compilation.Artifact, error) {
	// TODO this should be using "compilations", not singular compilation
	fReleaseQueryFilter, fReleaseQueryVarsTypes, fReleaseQueryVars := releaseversiongraphql.BuildListQueryArgs(f.Release)
	fOSQueryFilter, fOSQueryVarsTypes, fOSQueryVars := osversiongraphql.BuildListQueryArgs(f.OS)
	cmd := fmt.Sprintf(`query (%s, %s) {
  release (%s) {
		compilation (%s) {
			tarball {
				name
				size
				hashes {
					type
					hash
				}
				urls {
					url
				}
				metaurls {
					url
					mediatype
				}
			}
		}
	}
}`, fReleaseQueryVarsTypes, fOSQueryVarsTypes, fReleaseQueryFilter, fOSQueryFilter)

	req := graphql.NewRequest(cmd)

	for k, v := range fReleaseQueryVars {
		req.Var(k, v)
	}

	for k, v := range fOSQueryVars {
		req.Var(k, v)
	}

	var resp filterResponse

	err := i.client.Execute(req, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "executing remote request")
	}

	return []compilation.Artifact{resp.Release.Compilation}, nil
}

func (i *Index) StoreCompilationArtifact(_ compilation.Artifact) error {
	return datastore.UnsupportedOperationErr
}
