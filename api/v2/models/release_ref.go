package models

type ReleaseRef struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Checksum Checksum `json:"checksum"`
}
