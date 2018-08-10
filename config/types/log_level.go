package types

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type LogLevel logrus.Level

func (ll *LogLevel) UnmarshalYAML(data []byte) error {
	parsed, err := logrus.ParseLevel(string(data))
	if err != nil {
		return errors.Wrap(err, "parsing log level")
	}

	*ll = LogLevel(parsed)

	return nil
}
