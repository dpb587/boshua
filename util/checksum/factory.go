package checksum

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/dpb587/boshua/util/checksum/algorithm"
)

func New(a algorithm.Algorithm) *WritableChecksum {
	return &WritableChecksum{
		algorithm: a,
		hasher:    a.NewHash(),
	}
}

func NewExisting(a algorithm.Algorithm, d []byte) ImmutableChecksum {
	return ImmutableChecksum{
		algorithm: a,
		data:      d,
	}
}

func CreateFromString(raw string) (ImmutableChecksum, error) {
	split := strings.SplitN(raw, ":", 2)
	if len(split) == 1 {
		split = []string{"sha1", split[0]}
	}

	d, err := hex.DecodeString(split[1])
	if err != nil {
		return ImmutableChecksum{}, fmt.Errorf("decoding: %v", err)
	}

	a, err := algorithm.LookupName(split[0])
	if err == nil {
		a = algorithm.NewUnknown(split[0])
	}

	return NewExisting(a, d), nil
}
