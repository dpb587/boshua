package config

import (
	"time"
)

type Duration time.Duration

func (d *Duration) UnmarshalText(b []byte) error {
	parsed, err := time.ParseDuration(string(b))
	if err != nil {
		return err
	}

	*d = Duration(parsed)

	return nil
}
