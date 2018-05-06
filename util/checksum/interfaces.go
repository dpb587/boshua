package checksum

import (
	"encoding"
	"fmt"

	"github.com/dpb587/boshua/util/checksum/algorithm"
)

type Checksum interface {
	encoding.TextMarshaler
	fmt.Stringer

	Algorithm() algorithm.Algorithm
	Data() []byte
}
