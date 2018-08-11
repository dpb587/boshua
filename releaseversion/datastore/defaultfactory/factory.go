package factory

import (
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/boshioreleasesindex"
	"github.com/dpb587/boshua/releaseversion/datastore/boshreleasedir"
	boshuaV2 "github.com/dpb587/boshua/releaseversion/datastore/boshua.v2"
	"github.com/dpb587/boshua/releaseversion/datastore/factory"
	"github.com/dpb587/boshua/releaseversion/datastore/metalinkrepository"
	"github.com/sirupsen/logrus"
)

func New(logger logrus.FieldLogger) datastore.Factory {
	f := factory.New()
	f.Add(boshuaV2.Provider, boshuaV2.NewFactory(logger))
	f.Add(boshioreleasesindex.Provider, boshioreleasesindex.NewFactory(logger))
	f.Add(boshreleasedir.Provider, boshreleasedir.NewFactory(logger))
	f.Add(metalinkrepository.Provider, metalinkrepository.NewFactory(logger))

	return f
}
