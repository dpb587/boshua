package inmemory

import (
	"fmt"

	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/osversion/datastore"
)

type index struct {
	inmemory []osversion.Artifact

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

func (i *index) Find(ref osversion.Reference) (osversion.Artifact, error) {
	err := i.load()
	if err != nil {
		return osversion.Artifact{}, fmt.Errorf("reloading: %v", err)
	}

	for _, osversion := range i.inmemory {
		if osversion.Reference.Name != ref.Name {
			continue
		} else if osversion.Reference.Version != ref.Version {
			continue
		}

		return osversion, nil
	}

	return osversion.Artifact{}, datastore.MissingErr
}

func (i *index) List() ([]osversion.Artifact, error) {
	err := i.load()
	if err != nil {
		return []osversion.Artifact{}, fmt.Errorf("reloading: %v", err)
	}

	return i.inmemory, nil
}
