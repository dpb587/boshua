package output

import (
	"github.com/dpb587/boshua/checksum"
)

type Result struct {
	Path   string                       `json:"path" yaml:"path"`
	Result []checksum.ImmutableChecksum `json:"result" yaml:"result"`
}
