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

func (i *index) Find(ref compiledreleaseversion.Reference) (compiledreleaseversion.Artifact, error) {
	err := i.load()
	if err != nil {
		return compiledreleaseversion.Artifact{}, fmt.Errorf("reloading: %v", err)
	}

	for _, artifact := range i.inmemory {
		if artifact.ReleaseVersion.Name != ref.ReleaseVersion.Name {
			continue
		} else if artifact.ReleaseVersion.Version != ref.ReleaseVersion.Version {
			continue
		} else if artifact.StemcellVersion.OS != ref.StemcellVersion.OS {
			continue
		} else if artifact.StemcellVersion.Version != ref.StemcellVersion.Version {
			continue
		}

		for _, cs := range ref.ReleaseVersion.Checksums {
			if artifact.ReleaseVersion.Checksums.Contains(&cs) {
				return artifact, nil
			}
		}
	}

	return compiledreleaseversion.Artifact{}, datastore.MissingErr
}

func (i *index) List() ([]compiledreleaseversion.Artifact, error) {
	err := i.load()
	if err != nil {
		return []compiledreleaseversion.Artifact{}, fmt.Errorf("reloading: %v", err)
	}

	return i.inmemory, nil
}
