package manifest

import (
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Manifest struct {
	parsed       interface{}
	requirements []ReleasePatch
}

func (m *Manifest) Requirements() []ReleasePatch {
	return m.requirements
}

func (m *Manifest) UpdateRelease(release ReleasePatch) error {
	op := release.Op()

	updated, err := op.Apply(m.parsed)
	if err != nil {
		return errors.Wrap(err, "applying op")
	}

	m.parsed = updated

	return nil
}

func (m *Manifest) Bytes() ([]byte, error) {
	bytes, err := yaml.Marshal(m.parsed)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling yaml")
	}

	return bytes, nil
}
