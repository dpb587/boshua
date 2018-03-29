package releaseversions

type Checksum struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (c Checksum) Equals(cs Checksum) bool {
	if c.Type != cs.Type {
		return false
	} else if c.Value != cs.Value {
		return false
	}

	return true
}
