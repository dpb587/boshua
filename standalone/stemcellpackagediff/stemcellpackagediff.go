package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	analysisboshuaV2 "github.com/dpb587/boshua/analysis/datastore/boshua.v2"
	analysisscheduler "github.com/dpb587/boshua/analysis/datastore/scheduler"
	boshuaV2 "github.com/dpb587/boshua/artifact/datastore/datastoreutil/boshua.v2"
	"github.com/dpb587/boshua/metalink"
	stemcellpackagesV1 "github.com/dpb587/boshua/stemcellversion/analyzers/stemcellpackages.v1"
	stemcellpackagesV1result "github.com/dpb587/boshua/stemcellversion/analyzers/stemcellpackages.v1/result"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversionboshuaV2 "github.com/dpb587/boshua/stemcellversion/datastore/boshua.v2"
	"github.com/dpb587/boshua/task"
	schedulerboshuaV2 "github.com/dpb587/boshua/task/scheduler/boshua.v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	osArg, priorVersionArg, nextVersionArg, err := parseArgs(os.Args)
	if err != nil {
		panic(errors.Wrap(err, "parsing args"))
	}

	logger := newLogger()
	boshuaConfig := loadBoshuaConfig()

	index := stemcellversionboshuaV2.New(stemcellversionboshuaV2.Config{BoshuaConfig: boshuaConfig}, logger)
	analysisIndex := analysisscheduler.New(
		analysisboshuaV2.New(analysisboshuaV2.Config{BoshuaConfig: boshuaConfig}, logger),
		schedulerboshuaV2.New(schedulerboshuaV2.Config{BoshuaConfig: boshuaConfig}, logger),
		func(status task.Status) {
			fmt.Fprintf(os.Stderr, "%s [%s/%s] analysis is %s\n", time.Now().Format("15:04:05"), osArg, "something", status)
		},
	)

	packagesBefore, err := loadPackages(index, analysisIndex, osArg, priorVersionArg)
	if err != nil {
		panic(errors.Wrap(err, "loading before"))
	}

	packagesAfter, err := loadPackages(index, analysisIndex, osArg, nextVersionArg)
	if err != nil {
		panic(errors.Wrap(err, "loading after"))
	}

	packages := mergePackages(packagesBefore, packagesAfter)

	for _, pkg := range packages {
		if pkg.Before == nil {
			fmt.Printf("+ %s (%s)\n", pkg.Name, pkg.After.Version)
		} else if pkg.After == nil {
			fmt.Printf("- %s (%s)\n", pkg.Name, pkg.Before.Version)
		} else if pkg.Before.Version != pkg.After.Version {
			fmt.Printf("~ %s (%s --> %s)\n", pkg.Name, pkg.Before.Version, pkg.After.Version)
		}
	}
}

type PackageDiff struct {
	Name   string
	Before *stemcellpackagesV1result.RecordPackage
	After  *stemcellpackagesV1result.RecordPackage
}

func mergePackages(before, after []stemcellpackagesV1result.RecordPackage) []PackageDiff {
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

func parseArgs(args []string) (string, string, string, error) {
	if len(args) != 4 {
		return "", "", "", errors.New("expected 3 arguments")
	}

	return os.Args[1], os.Args[2], os.Args[3], nil
}

func newLogger() logrus.FieldLogger {
	var logger = logrus.New()

	logger.Out = os.Stderr
	logger.Formatter = &logrus.JSONFormatter{}

	if logLevel := os.Getenv("BOSHUA_LOG_LEVEL"); logLevel != "" {
		parsedLogLevel, err := logrus.ParseLevel(logLevel)
		if err != nil {
			panic(errors.Wrap(err, "parsing $BOSHUA_LOG_LEVEL"))
		}

		logger.Level = logrus.Level(parsedLogLevel)
	}

	return logger
}

func loadBoshuaConfig() boshuaV2.BoshuaConfig {
	server := os.Getenv("BOSHUA_SERVER")
	if server == "" {
		panic(errors.Wrap(errors.New("no boshua.v2 server specified"), "reading $BOSHUA_SERVER")) // TODO default to global
	}

	return boshuaV2.BoshuaConfig{
		URL: server,
	}
}

func loadPackages(index stemcellversiondatastore.Index, analysisIndex analysisdatastore.Index, os, version string) ([]stemcellpackagesV1result.RecordPackage, error) {
	ref := datastore.FilterParams{
		IaaSExpected:    true,
		IaaS:            "aws",
		FlavorExpected:  true,
		Flavor:          "light",
		OSExpected:      true,
		OS:              os,
		VersionExpected: true,
		Version:         version,
	}

	artifact, err := datastore.GetArtifact(index, ref)
	if err != nil {
		panic(errors.Wrap(err, "finding stemcell"))
	}

	analysis, err := analysisdatastore.GetAnalysisArtifact(analysisIndex, analysis.Reference{
		Analyzer: stemcellpackagesV1.AnalyzerName,
		Subject:  artifact,
	})
	if err != nil {
		return nil, errors.Wrap(err, "finding analysis")
	}

	r, w := io.Pipe()

	go func() {
		defer w.Close()

		err := metalink.StreamFile(analysis.MetalinkFile(), w)
		if err != nil {
			panic(errors.Wrap(err, "streaming"))
		}
	}()

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return nil, errors.Wrap(err, "starting gzip")
	}

	var pkgs []stemcellpackagesV1result.RecordPackage

	err = stemcellpackagesV1result.NewProcessor(gzr, func(r stemcellpackagesV1result.Record) error {
		if r.Package == nil {
			return nil
		}

		pkgs = append(pkgs, *r.Package)

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "processing old")
	}

	return pkgs, nil
}
