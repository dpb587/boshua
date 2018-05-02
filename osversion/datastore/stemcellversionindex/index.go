package boshio

import (
	"fmt"
	"reflect"

	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/osversion/datastore"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"

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

func (i *index) List() ([]osversion.Artifact, error) {
	matches := map[osversion.Reference]osversion.Artifact{}

	stemcells, err := i.stemcellVersionIndex.List()
	if err != nil {
		return nil, fmt.Errorf("listing stemcell versions: %v", err)
	}

	for _, stemcellVersion := range stemcells {
		if stemcellVersion.IaaS != "warden" {
			continue
		} else if stemcellVersion.Hypervisor != "boshlite" {
			continue
		}

		ref := osversion.Reference{
			Name:    stemcellVersion.OS,
			Version: stemcellVersion.Version,
		}

		matches[ref] = osversion.New(
			ref,
			stemcellVersion.MetalinkFile,
			stemcellVersion.MetalinkSource,
		)
	}

	var results []osversion.Artifact

	for _, artifact := range matches {
		results = append(results, artifact)
	}

	return results, nil
}

func (i *index) Find(ref osversion.Reference) (osversion.Artifact, error) {
	artifacts, err := i.List()
	if err != nil {
		return osversion.Artifact{}, fmt.Errorf("listing versions: %v", err)
	}

	for _, artifact := range artifacts {
		if artifact.Name != ref.Name {
			continue
		} else if artifact.Version != ref.Version {
			continue
		}

		return artifact, nil
	}

	return osversion.Artifact{}, datastore.MissingErr
}