package config

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

func RemarshalYAML(from interface{}, to interface{}) error {
	bytes, err := yaml.Marshal(from)
	if err != nil {
		return fmt.Errorf("marshaling: %v", err)
	}

	err = yaml.Unmarshal(bytes, to)
	if err != nil {
		return fmt.Errorf("unmarshalling: %v", err)
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
		return fmt.Errorf("marshaling: %v", err)
	}

	err = json.Unmarshal(bytes, to)
	if err != nil {
		return fmt.Errorf("unmarshalling: %v", err)
	}

	defaultable, ok := to.(Defaultable)
	if ok {
		defaultable.ApplyDefaults()
	}

	return nil
}
