package inmemory

import (
	"fmt"

	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
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
		return fmt.Errorf("checking reloader: %v", err)
	}

	if reload || i.inmemory == nil {
		return i.reload()
	}

	return nil
}

func (i *index) reload() error {
	data, err := i.loader()
	if err != nil {
		return fmt.Errorf("reloading: %v", err)
	}

	i.inmemory = data

	return nil
}

func (i *index) Filter(ref compiledreleaseversion.Reference) ([]compiledreleaseversion.Artifact, error) {
	var results []compiledreleaseversion.Artifact

	err := i.load()
	if err != nil {
		return nil, fmt.Errorf("reloading: %v", err)
	}

	for _, artifact := range i.inmemory {
		if artifact.ReleaseVersion.Name != ref.ReleaseVersion.Name {
			continue
		} else if artifact.ReleaseVersion.Version != ref.ReleaseVersion.Version {
			continue
		} else if artifact.OSVersion.Name != ref.OSVersion.Name {
			continue
		} else if artifact.OSVersion.Version != ref.OSVersion.Version {
			continue
		}

		for _, cs := range ref.ReleaseVersion.Checksums {
			if artifact.ReleaseVersion.Checksums.Contains(&cs) {
				results = append(results, artifact)

				break
			}
		}
	}

	return results, nil
}

func (i *index) Find(ref compiledreleaseversion.Reference) (compiledreleaseversion.Artifact, error) {
	return datastore.FilterForOne(i, ref)
}
