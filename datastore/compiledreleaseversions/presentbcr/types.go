package presentbcr

type Record struct {
	Release  RecordRelease  `json:"release"`
	Stemcell RecordStemcell `json:"stemcell"`
}

type RecordRelease struct {
	Name      string                  `json:"name"`
	Version   string                  `json:"version"`
	Checksums []RecordReleaseChecksum `json:"checksums"`
}

type RecordReleaseChecksum struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type RecordStemcell struct {
	OS      string `json:"os"`
	Version string `json:"version"`
}
