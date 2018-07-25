package formatter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/dpb587/boshua/analysis/analyzer/releasemanifests.v1/output"
	yaml "gopkg.in/yaml.v2"
)

type Properties struct {
	Jobs []string
}

func (f Properties) Format(writer io.Writer, reader io.Reader) error {
	s := bufio.NewScanner(reader)
	for s.Scan() {
		var result output.Result

		err := json.Unmarshal(s.Bytes(), &result)
		if err != nil {
			return err
		}

		var parsedManifest manifest

		err = yaml.Unmarshal([]byte(result.Raw), &parsedManifest)
		if err != nil {
			return err
		}

		if len(f.Jobs) > 0 {
			var found bool

			for _, job := range f.Jobs {
				if parsedManifest.Name == job {
					found = true

					break
				}
			}

			if !found {
				continue
			}
		}

		for propertyName, property := range parsedManifest.Properties {
			if len(f.Jobs) != 1 {
				fmt.Fprintf(writer, "%s\t", parsedManifest.Name)
			}

			fmt.Fprintf(writer, "%s\t%s\n", propertyName, strings.Replace(property.Description, `\n`, "\n", -1))
		}
	}

	return nil
}

type manifest struct {
	Name       string                      `yaml:"name"`
	Properties map[string]manifestProperty `yaml:"properties"`
}

type manifestProperty struct {
	Description string `yaml:"description"`
}
