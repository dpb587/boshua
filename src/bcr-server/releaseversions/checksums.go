package releaseversions

type Checksums []Checksum

func (c Checksums) Contains(checksum Checksum) bool {
	for _, cs := range c {
		if cs.Equals(checksum) {
			return true
		}
	}

	return false
}
