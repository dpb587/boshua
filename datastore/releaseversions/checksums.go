package releaseversions

type Checksums []Checksum

func (c Checksums) Preferred() Checksum {
	for _, algorithm := range []string{"sha256", "sha1"} {
		for _, checksum := range c {
			if checksum.Algorithm() == algorithm {
				return checksum
			}
		}
	}

	panic("missing sha1 or sha256")
}

func (c Checksums) Contains(expected Checksum) bool {
	for _, actual := range c {
		if actual == expected {
			return true
		}
	}

	return false
}
