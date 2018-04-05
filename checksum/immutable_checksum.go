package checksum

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/dpb587/bosh-compiled-releases/checksum/algorithm"
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
	var err error

	split := strings.SplitN(string(data), ":", 2)
	if len(split) == 1 {
		split = []string{"sha1", split[0]}
	}

	c.data, err = hex.DecodeString(split[1])
	if err != nil {
		return fmt.Errorf("decoding: %v", err)
	}

	c.algorithm, err = algorithm.LookupName(split[0])
	if err == nil {
		c.algorithm = algorithm.NewUnknown(split[0])
	}

	return nil
}

func (c *ImmutableChecksum) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%s:%x", c.algorithm.Name(), c.data)), nil
}

func (c *ImmutableChecksum) String() string {
	r, _ := c.MarshalText()
	return string(r)
}
