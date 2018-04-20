package inmemory

import (
	"fmt"

	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
)

type index struct {
	inmemory []compiledreleaseversion.Subject

	loader   Loader
	reloader Reloader

	releaseVersionIndex releaseversiondatastore.Index
}

var _ datastore.Index = &index{}

func New(loader Loader, reloader Reloader, releaseVersionIndex releaseversiondatastore.Index) datastore.Index {
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

func (i *index) Find(ref compiledreleaseversion.Reference) (compiledreleaseversion.Subject, error) {
	err := i.load()
	if err != nil {
		return compiledreleaseversion.Subject{}, fmt.Errorf("reloading: %v", err)
	}

	for _, subject := range i.inmemory {
		if subject.Release.Name != ref.Release.Name {
			continue
		} else if subject.Release.Version != ref.Release.Version {
			continue
		} else if subject.Stemcell.OS != ref.Stemcell.OS {
			continue
		} else if subject.Stemcell.Version != ref.Stemcell.Version {
			continue
		} else if subject.Release.Checksum.String() == ref.Release.Checksum.String() {
			// shortcut
			return subject, nil
		}

		// checksum matrix
		_, err := i.releaseVersionIndex.Find(subject.Release)
		if err == releaseversiondatastore.MissingErr {
			return compiledreleaseversion.Subject{}, datastore.MissingErr
		} else if err != nil {
			return compiledreleaseversion.Subject{}, fmt.Errorf("finding source: %v", err)
		}

		return subject, nil
	}

	return compiledreleaseversion.Subject{}, datastore.MissingErr
}

func (i *index) List() ([]compiledreleaseversion.Subject, error) {
	err := i.load()
	if err != nil {
		return []compiledreleaseversion.Subject{}, fmt.Errorf("reloading: %v", err)
	}

	return i.inmemory, nil
}
