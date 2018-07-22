package factory

import (
	"fmt"

	"github.com/dpb587/boshua/osversion/datastore"
	"github.com/dpb587/boshua/osversion/datastore/reflect"
	"github.com/dpb587/boshua/osversion/datastore/stemcellversionindex"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/sirupsen/logrus"
)

type factory struct {
	logger logrus.FieldLogger
}

var _ datastore.Factory = &factory{}

func New(logger logrus.FieldLogger) datastore.Factory {
	return &factory{
		logger: logger,
	}
}

func (f *factory) Create(provider, name string, options map[string]interface{}, stemcellVersionIndex stemcellversiondatastore.Index) (datastore.Index, error) {
	logger := f.logger.WithField("datastore", fmt.Sprintf("compiledreleaseversion:%s[%s]", provider, name))

	switch provider {
	case "stemcellversionindex":
		return stemcellversionindex.New(stemcellVersionIndex, logger), nil
	case "reflect":
		return reflect.New(logger), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
