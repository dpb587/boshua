package main

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	analysisV2 "github.com/dpb587/boshua/analysis/datastore/boshua.v2"
	"github.com/dpb587/boshua/cli/args"
	"github.com/dpb587/boshua/cli/opts"
	"github.com/dpb587/boshua/config"
	"github.com/dpb587/boshua/metalink/analysisprocessor"
	"github.com/dpb587/boshua/stemcellversion"
	stemcellpackagesV1result "github.com/dpb587/boshua/stemcellversion/analyzers/stemcellpackages.v1/result"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversionV2 "github.com/dpb587/boshua/stemcellversion/datastore/boshua.v2"
	schedulerV2 "github.com/dpb587/boshua/task/scheduler/boshua.v2"
	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type cmd struct {
	AppOpts *opts.Opts

	IaaS       string `long:"iaas" description:"Stemcell IaaS" default:"aws"`
	Hypervisor string `long:"hypervisor" description:"Stemcell hypervisor" default:"xen-hvm"`
	DiskFormat string `long:"disk-format" description:"Stemcell disk format"`
	Flavor     string `long:"flavor" description:"Stemcell flavor (heavy, light)" default:"heavy"`

	Format string `long:"format" description:"Output format (text, markdown, json)" default:"text"`
	Args   struct {
		OS     string `positional-arg-name:"OS" description:"Operating system name"`
		Before string `positional-arg-name:"BEFORE" description:"Earlier version"`
		After  string `positional-arg-name:"AFTER" description:"Later version"`
	} `positional-args:"true" required:"true"`
}

func (c *cmd) Execute(_ []string) error {
	cfg, err := c.AppOpts.GetConfig()
	if err != nil {
		return errors.Wrap(err, "loading config")
	}

	// only support remote api server
	cfg.SetAnalysisFactory(analysisV2.NewFactory(cfg.GetLogger()))
	cfg.SetStemcellFactory(stemcellversionV2.NewFactory(cfg.GetLogger()))
	cfg.SetSchedulerFactory(schedulerV2.NewFactory(cfg.GetLogger()))

	stemcellIndex, err := cfg.GetStemcellIndex(config.DefaultName)
	if err != nil {
		return errors.Wrap(err, "loading stemcell index")
	}

	analysisIndex, err := cfg.GetStemcellAnalysisIndex(config.DefaultName)
	if err != nil {
		return errors.Wrap(err, "loading analysis index")
	}

	refBefore := stemcellversion.Reference{
		IaaS:       c.IaaS,
		Hypervisor: c.Hypervisor,
		OS:         c.Args.OS,
		Version:    c.Args.Before,
		Flavor:     c.Flavor,
		DiskFormat: c.DiskFormat,
	}
	packagesBefore, err := loadPackages(stemcellIndex, analysisIndex, refBefore)
	if err != nil {
		return errors.Wrapf(err, "loading %s/%s", refBefore.FullName(), refBefore.Version)
	}

	refAfter := stemcellversion.Reference{
		IaaS:       c.IaaS,
		Hypervisor: c.Hypervisor,
		OS:         c.Args.OS,
		Version:    c.Args.After,
		Flavor:     c.Flavor,
		DiskFormat: c.DiskFormat,
	}
	packagesAfter, err := loadPackages(stemcellIndex, analysisIndex, refAfter)
	if err != nil {
		return errors.Wrapf(err, "loading %s/%s", refAfter.FullName(), c.Args.After)
	}

	var f Formatter

	switch c.Format {
	case "markdown":
		f = MarkdownFormatter{}
	case "text":
		f = TextFormatter{}
	case "json":
		f = JSONFormatter{}
	default:
		return fmt.Errorf("invalid format: %s", c.Format)
	}

	return f.Dump(os.Stdout, refBefore, refAfter, diffPackages(packagesBefore, packagesAfter))
}

var defaultServer string

func main() {
	c := cmd{
		AppOpts: &opts.Opts{
			DefaultServer: defaultServer,
			LogLevel:      args.LogLevel(logrus.FatalLevel),
		},
	}

	var parser = flags.NewParser(&c, flags.Default)
	parser.SubcommandsOptional = true
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		return c.Execute(args)
	}

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

func diffPackages(before, after []stemcellpackagesV1result.RecordPackage) []PackageDiff {
	mappedResults := map[string]PackageDiff{}

	for idx := range before {
		name := before[idx].Name

		mappedResults[name] = PackageDiff{
			Name:   name,
			Before: &before[idx],
		}
	}

	for idx := range after {
		name := after[idx].Name

		if _, found := mappedResults[name]; found {
			mappedResults[name] = PackageDiff{
				Name:   name,
				Before: mappedResults[name].Before,
				After:  &after[idx],
			}
		} else {
			mappedResults[name] = PackageDiff{
				Name:  name,
				After: &after[idx],
			}
		}
	}

	var results []PackageDiff

	for rIdx := range mappedResults {
		results = append(results, mappedResults[rIdx])
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results
}

func loadPackages(index stemcellversiondatastore.Index, analysisIndex analysisdatastore.Index, ref stemcellversion.Reference) ([]stemcellpackagesV1result.RecordPackage, error) {
	artifact, err := datastore.GetArtifact(index, datastore.FilterParamsFromReference(ref))
	if err != nil {
		return nil, errors.Wrap(err, "finding stemcell")
	}

	analysis, err := analysisdatastore.GetAnalysisArtifact(analysisIndex, analysis.Reference{
		Analyzer: "stemcellpackages.v1",
		Subject:  artifact,
	})
	if err != nil {
		return nil, errors.Wrap(err, "finding analysis")
	}

	var pkgs []stemcellpackagesV1result.RecordPackage

	err = analysisprocessor.Process(analysis, func(reader io.Reader) error {
		return stemcellpackagesV1result.NewProcessor(reader, func(r stemcellpackagesV1result.Record) error {
			if r.Package == nil {
				return nil
			}

			pkgs = append(pkgs, *r.Package)

			return nil
		})
	})
	if err != nil {
		return nil, errors.Wrap(err, "processing results")
	}

	return pkgs, nil
}
