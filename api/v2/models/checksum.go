package models

import "strings"

type Checksums []Checksum

type Checksum string

func (c Checksum) Algorithm() string {
	return c.tuple()[0]
}

func (c Checksum) Data() string {
	return c.tuple()[1]
}

func (c Checksum) tuple() [2]string {
	split := strings.SplitN(string(c), ":", 2)

	if len(split) == 2 {
		return [2]string{split[0], split[1]}
	}

	return [2]string{"sha1", split[0]}
}
