package inmemory

import (
	"fmt"

	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
)

type index struct {
	inmemory []releaseversion.Artifact

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

func (i *index) Find(ref releaseversion.Reference) (releaseversion.Artifact, error) {
	err := i.load()
	if err != nil {
		return releaseversion.Artifact{}, fmt.Errorf("reloading: %v", err)
	}

	for _, artifact := range i.inmemory {
		if artifact.Reference.Name != ref.Name {
			continue
		} else if artifact.Reference.Version != ref.Version {
			continue
		}

		if len(ref.Checksums) == 0 {
			return artifact, nil
		}

		for _, cs := range ref.Checksums.Prioritized() {
			if artifact.MatchesChecksum(&cs) {
				return artifact, nil
			}
		}
	}

	return releaseversion.Artifact{}, datastore.MissingErr
}

func (i *index) List() ([]releaseversion.Artifact, error) {
	err := i.load()
	if err != nil {
		return []releaseversion.Artifact{}, fmt.Errorf("reloading: %v", err)
	}

	return i.inmemory, nil
}
