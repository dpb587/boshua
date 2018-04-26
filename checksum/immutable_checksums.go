package checksum

type ImmutableChecksums []ImmutableChecksum

func (c ImmutableChecksums) Prioritized() ImmutableChecksums {
	var prioritized ImmutableChecksums

	for _, algorithm := range []string{"sha512", "sha256", "sha1", "md5"} {
		for _, checksum := range c {
			if checksum.Algorithm().Name() == algorithm {
				prioritized = append(prioritized, checksum)

				break
			}
		}
	}

	return prioritized
}

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
