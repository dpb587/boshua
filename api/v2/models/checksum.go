package models

type Checksums []Checksum

type Checksum struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
