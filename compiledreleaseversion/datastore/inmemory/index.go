package inmemory

import (
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/pkg/errors"
)

type index struct {
	inmemory []compiledreleaseversion.Artifact

	loader   Loader
	reloader Reloader
}

var _ datastore.Index = &index{}

func New(loader Loader, reloader Reloader) datastore.Index {
	return &index{
		loader:   loader,
		reloader: reloader,
	}
}

func (i *index) load() error {
	var reload bool
	var err error

	reload, err = i.reloader()
	if err != nil {
		return errors.Wrap(err, "checking reloader")
	}

	if reload || i.inmemory == nil {
		return i.reload()
	}

	return nil
}

func (i *index) reload() error {
	data, err := i.loader()
	if err != nil {
		return errors.Wrap(err, "reloading")
	}

	i.inmemory = data

	return nil
}

func (i *index) Filter(ref compiledreleaseversion.Reference) ([]compiledreleaseversion.Artifact, error) {
	var results []compiledreleaseversion.Artifact

	err := i.load()
	if err != nil {
		return nil, errors.Wrap(err, "reloading")
	}

	for _, artifact := range i.inmemory {
		artifactRef := artifact.Reference().(compiledreleaseversion.Reference)

		if artifactRef.ReleaseVersion.Name != ref.ReleaseVersion.Name {
			continue
		} else if artifactRef.ReleaseVersion.Version == "*" {
			// okay
		} else if artifactRef.ReleaseVersion.Version != ref.ReleaseVersion.Version {
			continue
		} else if artifactRef.OSVersion.Name != ref.OSVersion.Name {
			continue
		} else if artifactRef.OSVersion.Version == "*" {
			// okay
		} else if artifactRef.OSVersion.Version != ref.OSVersion.Version {
			continue
		}

		if len(ref.ReleaseVersion.Checksums) != 0 {
			var match int

			for _, cs := range ref.ReleaseVersion.Checksums {
				if artifactRef.ReleaseVersion.Checksums.Contains(&cs) {
					match += 1
				}
			}

			if match != len(ref.ReleaseVersion.Checksums) {
				continue
			}
		}

		results = append(results, artifact)
	}

	return results, nil
}

func (i *index) Find(ref compiledreleaseversion.Reference) (compiledreleaseversion.Artifact, error) {
	return datastore.FilterForOne(i, ref)
}

func (i *index) Store(artifact compiledreleaseversion.Artifact) error {
	return datastore.UnsupportedOperationErr
}

func (i *index) GetAnalysisDatastore(_ compiledreleaseversion.Reference) (analysisdatastore.Index, error) {
	return nil, datastore.UnsupportedOperationErr
}
