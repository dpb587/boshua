package args

import (
	"time"

	"github.com/pkg/errors"
)

type Duration time.Duration

func (d *Duration) UnmarshalFlag(data string) error {
	parsed, err := time.ParseDuration(data)
	if err != nil {
		return errors.Wrap(err, "parsing duration")
	}

	*d = Duration(parsed)

	return nil
}
