package inmemory

import (
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"fmt"
)

type index struct {
	inmemory []compiledreleaseversions.CompiledReleaseVersion

	loader              Loader
	releaseVersionIndex releaseversions.Index
}

func New(loader Loader, releaseVersionIndex releaseversions.Index) compiledreleaseversions.Index {
	return &index{
		loader:              loader,
		releaseVersionIndex: releaseVersionIndex,
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
		} else if compiledreleaseversion.Release.Checksum.Equals(ref.Release.Checksum) {
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
