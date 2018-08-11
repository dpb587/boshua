package factory

import (
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	boshuaV2 "github.com/dpb587/boshua/releaseversion/compilation/datastore/boshua.v2"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore/contextualosmetalinkrepository"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore/contextualrepoosmetalinkrepository"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore/factory"
	"github.com/sirupsen/logrus"
)

func New(logger logrus.FieldLogger) datastore.Factory {
	f := factory.New()
	f.Add(boshuaV2.Provider, boshuaV2.NewFactory(logger))
	f.Add(contextualosmetalinkrepository.Provider, contextualosmetalinkrepository.NewFactory(logger))
	f.Add(contextualrepoosmetalinkrepository.Provider, contextualrepoosmetalinkrepository.NewFactory(logger))

	return f
}
