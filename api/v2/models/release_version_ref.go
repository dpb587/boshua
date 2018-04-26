package models

import "github.com/dpb587/boshua/checksum"

type ReleaseVersionRef struct {
	Name     string                     `json:"name"`
	Version  string                     `json:"version"`
	Checksum checksum.ImmutableChecksum `json:"checksum"`
}
