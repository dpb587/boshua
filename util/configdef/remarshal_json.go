package configdef

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func RemarshalJSON(from interface{}, to interface{}) error {
	bytes, err := json.Marshal(from)
	if err != nil {
		return errors.Wrap(err, "marshaling")
	}

	return UnmarshalJSON(bytes, to)
}

func UnmarshalJSON(from []byte, to interface{}) error {
	err := json.Unmarshal(from, to)
	if err != nil {
		return errors.Wrap(err, "unmarshalling")
	}

	defaultable, ok := to.(Defaultable)
	if ok {
		defaultable.ApplyDefaults()
	}

	return nil
}
