package config

import (
	"encoding/json"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func RemarshalYAML(from interface{}, to interface{}) error {
	bytes, err := yaml.Marshal(from)
	if err != nil {
		return errors.Wrap(err, "marshaling")
	}

	return UnmarshalYAML(bytes, to)
}

func UnmarshalYAML(from []byte, to interface{}) error {
	err := yaml.Unmarshal(from, to)
	if err != nil {
		return errors.Wrap(err, "unmarshalling")
	}

	defaultable, ok := to.(Defaultable)
	if ok {
		defaultable.ApplyDefaults()
	}

	return nil
}

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
