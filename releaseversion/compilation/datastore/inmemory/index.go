package inmemory

import (
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
)

type Index struct {
	artifacts []compilation.Artifact
}

var _ datastore.Index = &Index{}

func New() *Index {
	return &Index{}
}

func (Index) GetName() string {
	panic("not supported directly")
}

func (i *Index) Add(artifact compilation.Artifact) {
	i.artifacts = append(i.artifacts, artifact)
}

func (i *Index) GetCompilationArtifacts(f datastore.FilterParams) ([]compilation.Artifact, error) {
	var results = []compilation.Artifact{}

	for _, artifact := range i.artifacts {
		if !f.OS.NameSatisfied(artifact.OS.Name) {
			continue
		} else if !f.OS.VersionSatisfied(artifact.OS.Version) {
			continue
		} else if !f.Release.NameSatisfied(artifact.Release.Name) {
			continue
		} else if !f.Release.VersionSatisfied(artifact.Release.Version) {
			continue
		} else if !f.Release.ChecksumSatisfied(metalinkutil.ChecksumsToHashes(artifact.Release.Checksums)) { // TODO avoid conversion
			continue
		}

		results = append(results, artifact)
	}

	return results, nil
}

func (i *Index) FlushCompilationCache() error {
	i.artifacts = []compilation.Artifact{}

	return nil
}

func (i *Index) StoreCompilationArtifact(artifact compilation.Artifact) error {
	return datastore.UnsupportedOperationErr
}
