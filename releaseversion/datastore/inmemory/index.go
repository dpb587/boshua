package inmemory

import (
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/pkg/errors"
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

func (i *index) Find(ref releaseversion.Reference) (releaseversion.Artifact, error) {
	return datastore.FilterForOne(i, ref)
}

func (i *index) Filter(ref releaseversion.Reference) ([]releaseversion.Artifact, error) {
	err := i.load()
	if err != nil {
		return nil, errors.Wrap(err, "reloading")
	}

	var results []releaseversion.Artifact

	for _, artifact := range i.inmemory {
		if artifact.Reference.Name != ref.Name {
			continue
		}

		if ref.Version == "*" {
			// okay
		} else if artifact.Reference.Version != ref.Version {
			continue
		}

		if len(ref.Checksums) != 0 {
			var match int

			for _, cs := range ref.Checksums.Prioritized() {
				if artifact.MatchesChecksum(&cs) {
					match += 1
				}
			}

			if match != len(ref.Checksums) {
				continue
			}
		}

		results = append(results, artifact)
	}

	return results, nil
}
