package formatter

import (
	"fmt"
	"io"

	"github.com/dpb587/boshua/pivnetfile/analyzers/tilereleasemanifests.v1/result"
	yaml "gopkg.in/yaml.v2"
)

type Releases struct {}

func (f Releases) Format(writer io.Writer, reader io.Reader) error {
	return result.NewProcessor(reader, func(record result.Record) error {
		if record.Path != "release.MF" {
			return nil
		}

		var parsedManifest manifest

		err := yaml.Unmarshal([]byte(record.Raw), &parsedManifest)
		if err != nil {
			return err
		}

		fmt.Fprintf(writer, "%s/%s", parsedManifest.Name, parsedManifest.Version)

		if len(parsedManifest.CompiledPackages) > 0 {
			// TODO verify all stemcells match?
			fmt.Fprintf(writer, "\t%s", parsedManifest.CompiledPackages[0].Stemcell)
		}

		fmt.Fprintf(writer, "\n")

		return nil
	})
}

type manifest struct {
	Name             string                    `yaml:"name"`
	Version          string                    `yaml:"version"`
	CompiledPackages []manifestCompiledPackage `yaml:"compiled_packages"`
}

type manifestCompiledPackage struct {
	Stemcell string `yaml:"stemcell"`
}
