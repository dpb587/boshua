package inmemory

import (
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
)

type Index struct {
	artifacts []releaseversion.Artifact
}

var _ datastore.Index = &Index{}

func New() *Index {
	return &Index{}
}

func (i *Index) Add(artifact releaseversion.Artifact) {
	i.artifacts = append(i.artifacts, artifact)
}

func (i *Index) GetArtifacts(f datastore.FilterParams) ([]releaseversion.Artifact, error) {
	var results = []releaseversion.Artifact{}

	for _, artifact := range i.artifacts {
		if !f.NameSatisfied(artifact.Name) {
			continue
		} else if !f.VersionSatisfied(artifact.Version) {
			continue
		} else if !f.ChecksumSatisfied(artifact.MetalinkFile().Hashes) {
			continue
		} else if !f.LabelsSatisfied(artifact.Labels) {
			continue
		}

		results = append(results, artifact)
	}

	return results, nil
}

func (i *Index) GetLabels() ([]string, error) {
	labelsMap := map[string]struct{}{}

	for _, one := range i.artifacts {
		for _, label := range one.Labels {
			labelsMap[label] = struct{}{}
		}
	}

	var labels []string

	for label := range labelsMap {
		labels = append(labels, label)
	}

	return labels, nil
}

func (i *Index) FlushCache() error {
	i.artifacts = []releaseversion.Artifact{}

	return nil
}
