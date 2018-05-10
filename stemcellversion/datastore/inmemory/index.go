package inmemory

import (
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/pkg/errors"
)

type index struct {
	inmemory []stemcellversion.Artifact

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
		return errors.Wrap(err, "checking reloader")
	}

	if reload || i.inmemory == nil {
		return i.reload()
	}

	return nil
}

func (i *index) reload() error {
	data, err := i.loader()
	if err != nil {
		return errors.Wrap(err, "reloading")
	}

	i.inmemory = data

	return nil
}

func (i *index) Filter(ref stemcellversion.Reference) ([]stemcellversion.Artifact, error) {
	var results []stemcellversion.Artifact

	err := i.load()
	if err != nil {
		return nil, errors.Wrap(err, "reloading")
	}

	for _, stemcellversion := range i.inmemory {
		if stemcellversion.Reference.IaaS != ref.IaaS {
			continue
		} else if stemcellversion.Reference.Hypervisor != ref.Hypervisor {
			continue
		} else if stemcellversion.Reference.OS != ref.OS {
			continue
		} else if stemcellversion.Reference.Light != ref.Light {
			continue
		}

		if ref.Version == "*" {
			// okay
		} else if stemcellversion.Reference.Version != ref.Version {
			continue
		}

		results = append(results, stemcellversion)
	}

	return results, nil
}

func (i *index) Find(ref stemcellversion.Reference) (stemcellversion.Artifact, error) {
	return datastore.FilterForOne(i, ref)
}

func (i *index) List() ([]stemcellversion.Artifact, error) {
	err := i.load()
	if err != nil {
		return nil, errors.Wrap(err, "reloading")
	}

	return i.inmemory, nil
}
