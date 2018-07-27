package boshuaV2

import (
	"fmt"
	"net/http"
	"reflect"

	boshuaV2 "github.com/dpb587/boshua/datastore/boshua.v2"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	datastoregraphql "github.com/dpb587/boshua/releaseversion/graphql"
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

func (i *Index) Filter(f *datastore.FilterParams) ([]releaseversion.Artifact, error) {
	fQueryFilter, fQueryVarsTypes, fQueryVars := datastoregraphql.BuildListQueryArgs(f)
	cmd := fmt.Sprintf(`query (%s) {
  releases (%s) {
    name
		version
		labels
		source_tarball {
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
}`, fQueryVarsTypes, fQueryFilter)

	req := graphql.NewRequest(cmd)

	for k, v := range fQueryVars {
		req.Var(k, v)
	}

	var resp filterResponse

	err := i.client.Execute(req, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "executing remote request")
	}

	return resp.Releases, nil
}

func (i *Index) Labels() ([]string, error) {
	return nil, errors.New("TODO")
}
