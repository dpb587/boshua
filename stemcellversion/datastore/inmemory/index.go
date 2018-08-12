package inmemory

import (
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
)

type Index struct {
	artifacts []stemcellversion.Artifact
}

var _ datastore.Index = &Index{}

func New() *Index {
	return &Index{}
}

func (i *Index) Add(artifact stemcellversion.Artifact) {
	i.artifacts = append(i.artifacts, artifact)
}

func (i *Index) GetArtifacts(f datastore.FilterParams) ([]stemcellversion.Artifact, error) {
	var results = []stemcellversion.Artifact{}

	for _, artifact := range i.artifacts {
		if !f.OSSatisfied(artifact.OS) {
			continue
		} else if !f.VersionSatisfied(artifact.Version) {
			continue
		} else if !f.IaaSSatisfied(artifact.IaaS) {
			continue
		} else if !f.HypervisorSatisfied(artifact.Hypervisor) {
			continue
		} else if !f.FlavorSatisfied(artifact.Flavor) {
			continue
		} else if !f.LabelsSatisfied(artifact.Labels) {
			return nil, nil
		}

		results = append(results, artifact)
	}

	return results, nil
}

func (i *Index) FlushCache() error {
	i.artifacts = []stemcellversion.Artifact{}

	return nil
}
