package boshuaV2

import (
	"fmt"
	"net/http"
	"reflect"

	boshuaV2 "github.com/dpb587/boshua/artifact/datastore/datastoreutil/boshua.v2"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	datastoregraphql "github.com/dpb587/boshua/releaseversion/graphql"
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

func (i *index) GetArtifacts(f datastore.FilterParams) ([]releaseversion.Artifact, error) {
	fQueryFilter, fQueryVarsTypes, fQueryVars := datastoregraphql.BuildListQueryArgs(f)
	if len(fQueryVarsTypes) > 0 {
		fQueryVarsTypes = fmt.Sprintf(`(%s)`, fQueryVarsTypes)
	}

	if len(fQueryFilter) > 0 {
		fQueryFilter = fmt.Sprintf(`(%s)`, fQueryFilter)
	}

	cmd := fmt.Sprintf(`query %s {
  releases %s {
    name
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

	results := resp.Releases

	for resultIdx := range results {
		results[resultIdx].Datastore = i.name
	}

	return results, nil
}

func (i *index) GetLabels() ([]string, error) {
	return nil, errors.New("TODO")
}

func (i *index) FlushCache() error {
	return nil
}
