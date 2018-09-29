package result

import (
	"bufio"
	"encoding/json"
	"io"
)

func NewProcessor(reader io.Reader, callback func(Record) error) error {
	s := bufio.NewScanner(reader)
	for s.Scan() {
		var r Record

		err := json.Unmarshal(s.Bytes(), &r)
		if err != nil {
			return err
		}

		err = callback(r)
		if err != nil {
			return err
		}
	}

	return nil
}
