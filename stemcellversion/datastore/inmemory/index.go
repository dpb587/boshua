package inmemory

import (
	"fmt"

	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
)

type index struct {
	inmemory []stemcellversion.Subject

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

func (i *index) Find(ref stemcellversion.Reference) (stemcellversion.Subject, error) {
	err := i.load()
	if err != nil {
		return stemcellversion.Subject{}, fmt.Errorf("reloading: %v", err)
	}

	for _, stemcellversion := range i.inmemory {
		if stemcellversion.Reference.OS != ref.OS {
			continue
		} else if stemcellversion.Reference.Version != ref.Version {
			continue
		}

		return stemcellversion, nil
	}

	return stemcellversion.Subject{}, datastore.MissingErr
}

func (i *index) List() ([]stemcellversion.Subject, error) {
	err := i.load()
	if err != nil {
		return []stemcellversion.Subject{}, fmt.Errorf("reloading: %v", err)
	}

	return i.inmemory, nil
}
