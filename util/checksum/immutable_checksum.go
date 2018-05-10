package checksum

import (
	"encoding/json"
	"fmt"

	"github.com/dpb587/boshua/util/checksum/algorithm"
	"github.com/pkg/errors"
)

type ImmutableChecksum struct {
	algorithm algorithm.Algorithm
	data      []byte
}

var _ Checksum = &ImmutableChecksum{}

func (c *ImmutableChecksum) Algorithm() algorithm.Algorithm {
	return c.algorithm
}

func (c *ImmutableChecksum) Data() []byte {
	return c.data
}

func (c *ImmutableChecksum) UnmarshalJSON(data []byte) error {
	var raw string

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	nc, err := CreateFromString(raw)
	if err != nil {
		return errors.Wrap(err, "parsing checksum")
	}

	*c = nc

	return nil
}

func (c ImmutableChecksum) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%s:%x", c.algorithm.Name(), c.data)), nil
}

func (c ImmutableChecksum) String() string {
	r, _ := c.MarshalText()
	return string(r)
}
