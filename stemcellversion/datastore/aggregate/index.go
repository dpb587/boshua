package aggregate

import (
	"fmt"

	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
)

type index struct {
	name    string
	indices []datastore.Index
}

var _ datastore.Index = &index{}

func New(name string, indices ...datastore.Index) datastore.Index {
	return &index{
		name:    name,
		indices: indices,
	}
}

func (i *index) GetName() string {
	return i.name
}

func (i *index) GetArtifacts(f datastore.FilterParams) ([]stemcellversion.Artifact, error) {
	// TODO merging behavior
	var results []stemcellversion.Artifact

	for idxIdx, idx := range i.indices {
		found, err := idx.GetArtifacts(f)
		if err != nil {
			return nil, fmt.Errorf("filtering %d: %v", idxIdx, err)
		}

		results = append(results, found...)
	}

	return results, nil
}

func (i *index) FlushCache() error {
	for idxIdx, idx := range i.indices {
		err := idx.FlushCache()
		if err != nil {
			return fmt.Errorf("flushing %d: %v", idxIdx, err)
		}
	}

	return nil
}
