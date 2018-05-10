package args

import (
	"github.com/dpb587/boshua/util/checksum"
	"github.com/pkg/errors"
)

type Checksum struct {
	checksum.ImmutableChecksum
}

func (c *Checksum) UnmarshalFlag(data string) error {
	nc, err := checksum.CreateFromString(data)
	if err != nil {
		return errors.Wrap(err, "parsing checksum arg")
	}

	c.ImmutableChecksum = nc

	return nil
}
