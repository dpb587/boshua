package factory

import (
	"github.com/dpb587/boshua/analysis/datastore"
	boshuaV2 "github.com/dpb587/boshua/analysis/datastore/boshua.v2"
	"github.com/dpb587/boshua/analysis/datastore/dpbreleaseartifacts"
	"github.com/dpb587/boshua/analysis/datastore/factory"
	"github.com/dpb587/boshua/analysis/datastore/localcache"
	"github.com/sirupsen/logrus"
)

func New(logger logrus.FieldLogger) datastore.Factory {
	f := factory.New()
	f.Add(boshuaV2.Provider, boshuaV2.NewFactory(logger))
	f.Add(dpbreleaseartifacts.Provider, dpbreleaseartifacts.NewFactory(logger))
	f.Add(localcache.Provider, localcache.NewFactory(logger))

	return f
}
