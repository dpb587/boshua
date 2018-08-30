package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/dpb587/boshua/stemcellversion"
	stemcellpackagesV1result "github.com/dpb587/boshua/stemcellversion/analyzers/stemcellpackages.v1/result"
)

type PackageDiff struct {
	Name   string                                  `json:"name" yaml:"name"`
	Before *stemcellpackagesV1result.RecordPackage `json:"before" yaml:"before"`
	After  *stemcellpackagesV1result.RecordPackage `json:"after" yaml:"after"`
}

type Formatter interface {
	Dump(io.Writer, stemcellversion.Reference, stemcellversion.Reference, []PackageDiff) error
}

type TextFormatter struct{}

func (f TextFormatter) Dump(w io.Writer, _ stemcellversion.Reference, _ stemcellversion.Reference, packages []PackageDiff) error {
	for _, pkg := range packages {
		if pkg.Before == nil {
			fmt.Fprintf(w, "+ %s (%s)\n", pkg.Name, pkg.After.Version)
		} else if pkg.After == nil {
			fmt.Fprintf(w, "- %s (%s)\n", pkg.Name, pkg.Before.Version)
		} else if pkg.Before.Version != pkg.After.Version {
			fmt.Fprintf(w, "~ %s (%s; was %s)\n", pkg.Name, pkg.After.Version, pkg.Before.Version)
		}
	}

	return nil
}

type JSONFormatter struct{}

func (f JSONFormatter) Dump(w io.Writer, _ stemcellversion.Reference, _ stemcellversion.Reference, packages []PackageDiff) error {
	var r []PackageDiff

	for _, pkg := range packages {
		if (pkg.Before == nil) || (pkg.After == nil) || (pkg.Before.Version != pkg.After.Version) {
			r = append(r, pkg)
		}
	}

	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintf(w, fmt.Sprintf("%s\n", b))

	return nil
}

type MarkdownFormatter struct{}

func (f MarkdownFormatter) Dump(w io.Writer, before stemcellversion.Reference, after stemcellversion.Reference, packages []PackageDiff) error {
	fmt.Fprintf(w, "| %s | %s | %s |\n", before.OS, before.Version, after.Version)
	fmt.Fprintf(w, "|:%s | %s:| %s:|\n", strings.Repeat("-", len(before.OS)), strings.Repeat("-", len(before.Version)), strings.Repeat("-", len(after.Version)))

	for _, pkg := range packages {
		if pkg.Before == nil {
			fmt.Fprintf(w, "| %s | &ndash; | %s |\n", pkg.Name, pkg.After.Version)
		} else if pkg.After == nil {
			fmt.Fprintf(w, "| %s | %s | &ndash; |\n", pkg.Name, pkg.Before.Version)
		} else if pkg.Before.Version != pkg.After.Version {
			fmt.Fprintf(w, "| %s | %s | %s |\n", pkg.Name, pkg.Before.Version, pkg.After.Version)
		}
	}

	return nil
}
