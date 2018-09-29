package factory

import (
	"github.com/dpb587/boshua/pivnetfile/datastore"
	"github.com/dpb587/boshua/pivnetfile/datastore/pivnet"
	"github.com/dpb587/boshua/pivnetfile/datastore/factory"
	"github.com/sirupsen/logrus"
)

func New(logger logrus.FieldLogger) datastore.Factory {
	f := factory.New()
	f.Add(pivnet.ProviderName, pivnet.NewFactory(logger))

	return f
}
