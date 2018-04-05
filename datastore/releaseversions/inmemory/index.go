package inmemory

import (
	"fmt"

	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
)

type index struct {
	inmemory []releaseversions.ReleaseVersion

	loader   Loader
	reloader Reloader
}

func New(loader Loader, reloader Reloader) releaseversions.Index {
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

func (i *index) Find(ref releaseversions.ReleaseVersionRef) (releaseversions.ReleaseVersion, error) {
	err := i.load()
	if err != nil {
		return releaseversions.ReleaseVersion{}, fmt.Errorf("reloading: %v", err)
	}

	for _, releaseversion := range i.inmemory {
		if releaseversion.ReleaseVersionRef.Name != ref.Name {
			continue
		} else if releaseversion.ReleaseVersionRef.Version != ref.Version {
			continue
		}

		if releaseversion.Checksums.Contains(ref.Checksum) {
			return releaseversion, nil
		}
	}

	return releaseversions.ReleaseVersion{}, releaseversions.MissingErr
}

func (i *index) List() ([]releaseversions.ReleaseVersion, error) {
	err := i.load()
	if err != nil {
		return []releaseversions.ReleaseVersion{}, fmt.Errorf("reloading: %v", err)
	}

	return i.inmemory, nil
}
