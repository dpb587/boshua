package checksum

import (
	"encoding"
	"fmt"

	"github.com/dpb587/boshua/checksum/algorithm"
)

type Checksum interface {
	encoding.TextMarshaler
	fmt.Stringer

	Algorithm() algorithm.Algorithm
	Data() []byte
}
