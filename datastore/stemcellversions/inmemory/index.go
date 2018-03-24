package inmemory

import (
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
	"fmt"
)

type index struct {
	inmemory []stemcellversions.StemcellVersion

	loader Loader
}

func New(loader Loader) stemcellversions.Index {
	return &index{
		loader: loader,
	}
}

func (i *index) load() error {
	if i.inmemory == nil {
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
