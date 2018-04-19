package inmemory

import (
	"fmt"

	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
)

type index struct {
	inmemory []stemcellversions.StemcellVersion

	loader   Loader
	reloader Reloader
}

var _ stemcellversions.Index = &index{}

func New(loader Loader, reloader Reloader) stemcellversions.Index {
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

func (i *index) Find(ref stemcellversions.StemcellVersionRef) (stemcellversions.StemcellVersion, error) {
	err := i.load()
	if err != nil {
		return stemcellversions.StemcellVersion{}, fmt.Errorf("reloading: %v", err)
	}

	for _, stemcellversion := range i.inmemory {
		if stemcellversion.StemcellVersionRef.OS != ref.OS {
			continue
		} else if stemcellversion.StemcellVersionRef.Version != ref.Version {
			continue
		}

		return stemcellversion, nil
	}

	return stemcellversions.StemcellVersion{}, stemcellversions.MissingErr
}

func (i *index) List() ([]stemcellversions.StemcellVersion, error) {
	err := i.load()
	if err != nil {
		return []stemcellversions.StemcellVersion{}, fmt.Errorf("reloading: %v", err)
	}

	return i.inmemory, nil
}
