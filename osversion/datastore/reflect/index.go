package reflect

import (
	"reflect"

	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/osversion/datastore"
	"github.com/dpb587/metalink"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger logrus.FieldLogger
}

var _ datastore.Index = &index{}

func New(logger logrus.FieldLogger) datastore.Index {
	return &index{
		logger: logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
	}
}

func (i *index) Filter(ref osversion.Reference) ([]osversion.Artifact, error) {
	return []osversion.Artifact{osversion.New(ref, metalink.File{})}, nil
}

func (i *index) Find(ref osversion.Reference) (osversion.Artifact, error) {
	return datastore.FilterForOne(i, ref)
}
