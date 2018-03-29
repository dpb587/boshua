package releaseversions

type Checksums []Checksum

func (c Checksums) Preferred() string {
	for _, t := range []string{"sha256", "sha1"} {
		for _, checksum := range c {
			if checksum.Type == t {
				return checksum.Value
			}
		}
	}

	panic("missing sha1 or sha256")
}

func (c Checksums) Contains(checksum Checksum) bool {
	for _, cs := range c {
		if cs.Equals(checksum) {
			return true
		}
	}

	return false
}
