package analyzer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/releaseversion/analyzers/releasemanifests.v1/output"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type analysisGenerator struct {
	tarball string
}

var _ analysis.AnalysisGenerator = &analysisGenerator{}

func NewAnalysis(tarball string) analysis.AnalysisGenerator {
	return &analysisGenerator{
		tarball: tarball,
	}
}

func (a *analysisGenerator) Analyze(results analysis.Writer) error {
	fh, err := os.Open(a.tarball)
	if err != nil {
		return errors.Wrap(err, "opening file")
	}

	defer fh.Close()

	gzReader, err := gzip.NewReader(fh)
	if err != nil {
		return errors.Wrap(err, "starting gzip")
	}

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrap(err, "advancing tar")
		} else if header.Typeflag == tar.TypeDir {
			continue
		}

		path := strings.TrimPrefix(header.Name, "./")

		if path == "release.MF" {
			err = a.analyzeReleaseManifest(results, path, tarReader)
		} else if strings.HasPrefix(path, "jobs/") && strings.HasSuffix(path, ".tgz") {
			err = a.analyzeJobArtifactManifest(results, path, tarReader)
		} else {
			continue
		}

		if err != nil {
			return fmt.Errorf("analyzing artifact %s: %v", path, err)
		}
	}

	return nil
}

func (a *analysisGenerator) analyzeReleaseManifest(results analysis.Writer, artifact string, reader io.Reader) error {
	marshalBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.Wrap(err, "reading release.MF")
	}

	var spec output.ResultSpec

	err = yaml.Unmarshal(marshalBytes, &spec)
	if err != nil {
		return errors.Wrap(err, "parsing release.MF")
	}

	err = results.Write(output.Result{
		Path:   artifact,
		Raw:    string(marshalBytes),
		Parsed: safejson(spec).(output.ResultSpec),
	})
	if err != nil {
		return errors.Wrap(err, "writing result")
	}

	return nil
}

func (a *analysisGenerator) analyzeJobArtifactManifest(results analysis.Writer, artifact string, reader io.Reader) error {
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return errors.Wrap(err, "starting gzip")
	}

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrap(err, "advancing tar")
		} else if header.Typeflag == tar.TypeDir {
			continue
		}

		path := strings.TrimPrefix(header.Name, "./")

		if path != "job.MF" {
			continue
		}

		marshalBytes, err := ioutil.ReadAll(tarReader)
		if err != nil {
			return errors.Wrap(err, "reading job.MF")
		}

		var spec output.ResultSpec

		err = yaml.Unmarshal(marshalBytes, &spec)
		if err != nil {
			return errors.Wrap(err, "parsing job.MF")
		}

		err = results.Write(output.Result{
			Path:   artifact,
			Raw:    string(marshalBytes),
			Parsed: safejson(spec).(output.ResultSpec),
		})
		if err != nil {
			return errors.Wrap(err, "writing result")
		}
	}

	return nil
}

func safejson(v interface{}) interface{} {
	switch tv := v.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for tvk, tvv := range tv {
			m2[tvk.(string)] = safejson(tvv)
		}
		return m2
	case []interface{}:
		for tvk, tvv := range tv {
			tv[tvk] = safejson(tvv)
		}
	}

	return v
}
