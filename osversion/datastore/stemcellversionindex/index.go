package boshio

import (
	"reflect"

	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/osversion/datastore"
	"github.com/dpb587/boshua/stemcellversion"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger               logrus.FieldLogger
	stemcellVersionIndex stemcellversiondatastore.Index
}

var _ datastore.Index = &index{}

func New(stemcellVersionIndex stemcellversiondatastore.Index, logger logrus.FieldLogger) datastore.Index {
	return &index{
		logger:               logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		stemcellVersionIndex: stemcellVersionIndex,
	}
}

func (i *index) Filter(ref osversion.Reference) ([]osversion.Artifact, error) {
	artifacts, err := i.list()
	if err != nil {
		return nil, errors.Wrap(err, "listing versions")
	}

	var results []osversion.Artifact

	for _, artifact := range artifacts {
		artifactRef := artifact.Reference().(osversion.Reference)

		if artifactRef.Name != ref.Name {
			continue
		}

		if ref.Version == "*" {
			// okay
		} else if artifactRef.Version != ref.Version {
			continue
		}

		results = append(results, artifact)
	}

	return results, nil
}

func (i *index) Find(ref osversion.Reference) (osversion.Artifact, error) {
	return datastore.FilterForOne(i, ref)
}

func (i *index) list() ([]osversion.Artifact, error) {
	matches := map[osversion.Reference]osversion.Artifact{}

	stemcells, err := i.stemcellVersionIndex.List()
	if err != nil {
		return nil, errors.Wrap(err, "listing stemcell versions")
	}

	for _, artifact := range stemcells {
		artifactRef := artifact.Reference().(stemcellversion.Reference)

		if artifactRef.IaaS != "warden" {
			continue
		} else if artifactRef.Hypervisor != "boshlite" {
			continue
		}

		ref := osversion.Reference{
			Name:    artifactRef.OS,
			Version: artifactRef.Version,
		}

		matches[ref] = osversion.New(
			ref,
			artifact.MetalinkFile(),
		)
	}

	var results []osversion.Artifact

	for _, artifact := range matches {
		results = append(results, artifact)
	}

	return results, nil
}
