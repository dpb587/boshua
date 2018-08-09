package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/dpb587/boshua/releaseversion/analyzers/releasemanifests.v1/result"
	yaml "gopkg.in/yaml.v2"
)

type Properties struct {
	Jobs []string
}

func (f Properties) Format(writer io.Writer, reader io.Reader) error {
	return result.NewProcessor(reader, func(record result.Record) error {
		var parsedManifest manifest

		err := yaml.Unmarshal([]byte(record.Raw), &parsedManifest)
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
				return nil
			}
		}

		for propertyName, property := range parsedManifest.Properties {
			if len(f.Jobs) != 1 {
				fmt.Fprintf(writer, "%s\t", parsedManifest.Name)
			}

			fmt.Fprintf(writer, "%s\t%s\n", propertyName, strings.Replace(property.Description, `\n`, "\n", -1))
		}

		return nil
	})
}

type manifest struct {
	Name       string                      `yaml:"name"`
	Properties map[string]manifestProperty `yaml:"properties"`
}

type manifestProperty struct {
	Description string `yaml:"description"`
}
