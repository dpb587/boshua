package factory

import (
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	boshuaV2 "github.com/dpb587/boshua/releaseversion/compilation/datastore/boshua.v2"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore/contextualosmetalinkrepository"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore/contextualrepoosmetalinkrepository"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore/factory"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/sirupsen/logrus"
)

func New(releaseVersionGetter releaseversiondatastore.NamedGetter, logger logrus.FieldLogger) datastore.Factory {
	f := factory.New()
	f.Add(boshuaV2.ProviderName, boshuaV2.NewFactory(logger))
	f.Add(contextualosmetalinkrepository.ProviderName, contextualosmetalinkrepository.NewFactory(releaseVersionGetter, logger))
	f.Add(contextualrepoosmetalinkrepository.ProviderName, contextualrepoosmetalinkrepository.NewFactory(releaseVersionGetter, logger))

	return f
}
