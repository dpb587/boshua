package releaseversion

import "github.com/dpb587/boshua/checksum"

type Reference struct {
	Name     string                     `json:"name"`
	Version  string                     `json:"version"`
	Checksum checksum.ImmutableChecksum `json:"checksum"`
}
