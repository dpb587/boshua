package legacybcr

import "github.com/dpb587/boshua/checksum"

type Record struct {
	Name     string         `json:"name"`
	Version  string         `json:"version"`
	Source   RecordSource   `json:"source"`
	Stemcell RecordStemcell `json:"stemcell"`
	Tarball  RecordTarball  `json:"tarball"`
}

type RecordSource struct {
	Digest checksum.ImmutableChecksum `json:"digest"`
}

type RecordStemcell struct {
	OS      string `json:"os"`
	Version string `json:"version"`
}

type RecordTarball struct {
	Digest checksum.ImmutableChecksum `json:"digest"`
	URL    string                     `json:"url"`
}
