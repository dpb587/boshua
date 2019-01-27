package checksum

import "fmt"

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
	for _, name := range []string{"sha256", "sha1"} {
		checksum, err := c.GetByAlgorithm(name)
		if err == nil {
			return checksum
		}
	}

	panic("missing sha1 or sha256")
}

func (c ImmutableChecksums) GetByAlgorithm(name string) (ImmutableChecksum, error) {
	for _, checksum := range c {
		if checksum.Algorithm().Name() == name {
			return checksum, nil
		}
	}

	return ImmutableChecksum{}, fmt.Errorf("unable to find checksum: %s", name)
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
