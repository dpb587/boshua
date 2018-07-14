package boshuav2

import (
	"net/http"
	"reflect"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/datastore/boshua.v2"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger logrus.FieldLogger
	client *boshuav2.Client
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	return &index{
		logger: logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		client: boshuav2.NewClient(http.DefaultClient, config.BoshuaConfig, logger),
	}
}

func (i *index) Filter(ref releaseversion.Reference) ([]releaseversion.Artifact, error) {
	return i.inmemory.Filter(ref)
}

func (i *index) Find(ref releaseversion.Reference) (releaseversion.Artifact, error) {
	return i.inmemory.Find(ref)
}

func (i *index) GetAnalysisDatastore(_ releaseversion.Reference) (analysisdatastore.Index, error) {
	return nil, datastore.UnsupportedOperationErr
}
