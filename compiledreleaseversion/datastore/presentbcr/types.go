package presentbcr

import (
	"github.com/dpb587/boshua/checksum"
)

type Record struct {
	Release  RecordRelease  `json:"release"`
	Stemcell RecordStemcell `json:"stemcell"`
}

type RecordRelease struct {
	Name      string                       `json:"name"`
	Version   string                       `json:"version"`
	Checksums []checksum.ImmutableChecksum `json:"checksums"`
}

type RecordStemcell struct {
	OS      string `json:"os"`
	Version string `json:"version"`
}
