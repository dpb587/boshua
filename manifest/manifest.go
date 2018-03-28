package manifest

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type Manifest struct {
	parsed       interface{}
	requirements []Release
}

func (m *Manifest) Requirements() []Release {
	return m.requirements
}

func (m *Manifest) UpdateRelease(release Release) error {
	op := release.Op()

	updated, err := op.Apply(m.parsed)
	if err != nil {
		return fmt.Errorf("applying op: %v", err)
	}

	m.parsed = updated

	return nil
}

func (m *Manifest) Bytes() ([]byte, error) {
	bytes, err := yaml.Marshal(m.parsed)
	if err != nil {
		return nil, fmt.Errorf("marshalling yaml: %v", err)
	}

	return bytes, nil
}
