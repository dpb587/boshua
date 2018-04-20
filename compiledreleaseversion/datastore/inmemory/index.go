package inmemory

import (
	"fmt"

	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore"
)

type index struct {
	inmemory []compiledreleaseversions.CompiledReleaseVersion

	loader   Loader
	reloader Reloader

	releaseVersionIndex releaseversions.Index
}

var _ compiledreleaseversions.Index = &index{}

func New(loader Loader, reloader Reloader, releaseVersionIndex releaseversions.Index) compiledreleaseversions.Index {
	return &index{
		loader:              loader,
		reloader:            reloader,
		releaseVersionIndex: releaseVersionIndex,
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

func (i *index) Find(ref compiledreleaseversions.CompiledReleaseVersionRef) (compiledreleaseversions.CompiledReleaseVersion, error) {
	err := i.load()
	if err != nil {
		return compiledreleaseversions.CompiledReleaseVersion{}, fmt.Errorf("reloading: %v", err)
	}

	for _, compiledreleaseversion := range i.inmemory {
		if compiledreleaseversion.Release.Name != ref.Release.Name {
			continue
		} else if compiledreleaseversion.Release.Version != ref.Release.Version {
			continue
		} else if compiledreleaseversion.Stemcell.OS != ref.Stemcell.OS {
			continue
		} else if compiledreleaseversion.Stemcell.Version != ref.Stemcell.Version {
			continue
		} else if compiledreleaseversion.Release.Checksum.String() == ref.Release.Checksum.String() {
			// shortcut
			return compiledreleaseversion, nil
		}

		// checksum matrix
		_, err := i.releaseVersionIndex.Find(compiledreleaseversion.Release)
		if err == releaseversions.MissingErr {
			return compiledreleaseversions.CompiledReleaseVersion{}, compiledreleaseversions.MissingErr
		} else if err != nil {
			return compiledreleaseversions.CompiledReleaseVersion{}, fmt.Errorf("finding source: %v", err)
		}

		return compiledreleaseversion, nil
	}

	return compiledreleaseversions.CompiledReleaseVersion{}, compiledreleaseversions.MissingErr
}

func (i *index) List() ([]compiledreleaseversions.CompiledReleaseVersion, error) {
	err := i.load()
	if err != nil {
		return []compiledreleaseversions.CompiledReleaseVersion{}, fmt.Errorf("reloading: %v", err)
	}

	return i.inmemory, nil
}
