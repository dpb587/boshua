package boshuaV2

import (
	"fmt"
	"net/http"
	"reflect"

	boshuaV2 "github.com/dpb587/boshua/artifact/datastore/datastoreutil/boshua.v2"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	datastoregraphql "github.com/dpb587/boshua/stemcellversion/graphql"
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger logrus.FieldLogger
	config Config
	client *boshuaV2.Client
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		logger: logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config: config,
		client: boshuaV2.NewClient(http.DefaultClient, config.BoshuaConfig, logger),
	}
}

func (i *index) GetArtifacts(f datastore.FilterParams) ([]stemcellversion.Artifact, error) {
	fQueryFilter, fQueryVarsTypes, fQueryVars := datastoregraphql.BuildListQueryArgs(f)
	cmd := fmt.Sprintf(`query (%s) {
  stemcells (%s) {
    os
		version
		iaas
		hypervisor
		diskFormat
		flavor
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

	return resp.Stemcells, nil
}
