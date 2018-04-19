package checksum

import (
	"fmt"
	"hash"

	"github.com/dpb587/boshua/checksum/algorithm"
)

type WritableChecksum struct {
	algorithm algorithm.Algorithm
	hasher    hash.Hash
}

var _ Checksum = &WritableChecksum{}

func (c *WritableChecksum) Algorithm() algorithm.Algorithm {
	return c.algorithm
}

func (c *WritableChecksum) Data() []byte {
	return c.hasher.Sum(nil)
}

func (c *WritableChecksum) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%s:%x", c.algorithm.Name(), c.Data())), nil
}

func (c *WritableChecksum) Write(p []byte) (int, error) {
	return c.hasher.Write(p)
}

func (c *WritableChecksum) String() string {
	r, _ := c.MarshalText()
	return string(r)
}

func (c *WritableChecksum) ImmutableChecksum() ImmutableChecksum {
	return NewExisting(c.algorithm, c.Data())
}
