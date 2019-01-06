package boshuaV2

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	boshuaV2 "github.com/dpb587/boshua/artifact/datastore/datastoreutil/boshua.v2"
	"github.com/dpb587/boshua/osversion"
	osversiongraphql "github.com/dpb587/boshua/osversion/graphql"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	releaseversiongraphql "github.com/dpb587/boshua/releaseversion/graphql"
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type index struct {
	name   string
	logger logrus.FieldLogger
	config Config
	client *boshuaV2.Client
}

var _ datastore.Index = &index{}

func New(name string, config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		name:   name,
		logger: logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config: config,
		client: boshuaV2.NewClient(http.DefaultClient, config.BoshuaConfig, logger),
	}
}

func (i *index) GetName() string {
	return i.name
}

func (i *index) GetCompilationArtifacts(f datastore.FilterParams) ([]compilation.Artifact, error) {
	// TODO this should be using "compilations", not singular compilation
	fReleaseQueryFilter, fReleaseQueryVarsTypes, fReleaseQueryVars := releaseversiongraphql.BuildListQueryArgs(
		f.Release,
		releaseversiondatastore.SingleArtifactLimitParams,
	)
	if len(fReleaseQueryFilter) > 0 {
		fReleaseQueryFilter = fmt.Sprintf(`(%s)`, fReleaseQueryFilter)
	}

	fOSQueryFilter, fOSQueryVarsTypes, fOSQueryVars := osversiongraphql.BuildListQueryArgs(f.OS)
	if len(fOSQueryFilter) > 0 {
		fOSQueryFilter = fmt.Sprintf(`(%s)`, fOSQueryFilter)
	}

	fQueryVarsTypes := strings.Join([]string{fReleaseQueryVarsTypes, fOSQueryVarsTypes}, ", ")
	if len(fQueryVarsTypes) > 0 {
		fQueryVarsTypes = fmt.Sprintf(`(%s)`, fQueryVarsTypes)
	}

	// TODO weird singular vs multiple queries

	cmd := fmt.Sprintf(`query %s {
  releases %s {
		name
		version
		labels

    compilations %s {
			os
			version
			labels

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
}`, fQueryVarsTypes, fReleaseQueryFilter, fOSQueryFilter)

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

	var results []compilation.Artifact

	for _, respRelease := range resp.Releases {
		for _, compl := range respRelease.Compilations {
			results = append(results, compilation.Artifact{
				Datastore: i.name,
				OS:        osversion.Reference{Name: compl.OS, Version: compl.Version},
				Release: releaseversion.Reference{
					Name:    respRelease.Name,
					Version: respRelease.Version,
				},
				Tarball: compl.Tarball,
				Labels:  compl.Labels,
			})
		}
	}

	return results, nil
}

func (i *index) StoreCompilationArtifact(_ compilation.Artifact) error {
	return datastore.UnsupportedOperationErr
}

func (i *index) FlushCompilationCache() error {
	return nil
}
