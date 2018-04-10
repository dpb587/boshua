package args

import (
	"fmt"

	"github.com/dpb587/bosh-compiled-releases/checksum"
)

type Checksum struct {
	checksum.ImmutableChecksum
}

func (c *Checksum) UnmarshalFlag(data string) error {
	nc, err := checksum.CreateFromString(data)
	if err != nil {
		return fmt.Errorf("parsing checksum arg: %v", err)
	}

	c.ImmutableChecksum = nc

	return nil
}
