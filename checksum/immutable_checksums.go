package checksum

type ImmutableChecksums []ImmutableChecksum

func (c ImmutableChecksums) Preferred() ImmutableChecksum {
	for _, algorithm := range []string{"sha256", "sha1"} {
		for _, checksum := range c {
			if checksum.Algorithm().Name() == algorithm {
				return checksum
			}
		}
	}

	panic("missing sha1 or sha256")
}

func (c ImmutableChecksums) Contains(expected Checksum) bool {
	expectedString := expected.String()

	for _, actual := range c {
		if actual.String() == expectedString {
			return true
		}
	}

	return false
}
