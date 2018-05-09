package args

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type LogLevel logrus.Level

func (ll *LogLevel) UnmarshalFlag(data string) error {
	parsed, err := logrus.ParseLevel(data)
	if err != nil {
		return fmt.Errorf("parsing log level: %v", err)
	}

	*ll = LogLevel(parsed)

	return nil
}
