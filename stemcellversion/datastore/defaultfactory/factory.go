package factory

import (
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore/boshioindex"
	boshuaV2 "github.com/dpb587/boshua/stemcellversion/datastore/boshua.v2"
	"github.com/dpb587/boshua/stemcellversion/datastore/factory"
	"github.com/sirupsen/logrus"
)

func New(logger logrus.FieldLogger) datastore.Factory {
	f := factory.New()
	f.Add(boshuaV2.ProviderName, boshuaV2.NewFactory(logger))
	f.Add(boshioindex.ProviderName, boshioindex.NewFactory(logger))

	return f
}
