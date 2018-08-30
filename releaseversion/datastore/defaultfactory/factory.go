package factory

import (
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/boshioindex"
	"github.com/dpb587/boshua/releaseversion/datastore/boshreleasedir"
	boshuaV2 "github.com/dpb587/boshua/releaseversion/datastore/boshua.v2"
	"github.com/dpb587/boshua/releaseversion/datastore/factory"
	"github.com/dpb587/boshua/releaseversion/datastore/metalinkrepository"
	"github.com/dpb587/boshua/releaseversion/datastore/trustedtarball"
	"github.com/sirupsen/logrus"
)

func New(logger logrus.FieldLogger) datastore.Factory {
	f := factory.New()
	f.Add(boshuaV2.ProviderName, boshuaV2.NewFactory(logger))
	f.Add(boshioindex.ProviderName, boshioindex.NewFactory(logger))
	f.Add(boshreleasedir.ProviderName, boshreleasedir.NewFactory(logger))
	f.Add(metalinkrepository.ProviderName, metalinkrepository.NewFactory(logger))
	f.Add(trustedtarball.ProviderName, trustedtarball.NewFactory(logger))

	return f
}
