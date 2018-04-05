package checksum

import (
	"encoding"
	"fmt"

	"github.com/dpb587/bosh-compiled-releases/checksum/algorithm"
)

type Checksum interface {
	encoding.TextMarshaler
	fmt.Stringer

	Algorithm() algorithm.Algorithm
	Data() []byte
}
