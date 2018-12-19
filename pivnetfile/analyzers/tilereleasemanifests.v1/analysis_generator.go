package analyzer

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/pivnetfile/analyzers/tilereleasemanifests.v1/result"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type analysisGenerator struct {
	subject string
}

var _ analysis.AnalysisGenerator = &analysisGenerator{}

func NewAnalysis(subject string) analysis.AnalysisGenerator {
	return &analysisGenerator{
		subject: subject,
	}
}

func (a *analysisGenerator) Analyze(records analysis.Writer) error {
	fh, err := zip.OpenReader(a.subject)
	if err != nil {
		return errors.Wrap(err, "opening file")
	}

	defer fh.Close()

	for _, f := range fh.File {
		path := strings.TrimPrefix(f.Name, "./")

		// TODO swith to something less generic and fragile
		if !strings.HasPrefix(path, "releases/") {
			continue
		} else if !strings.HasSuffix(path, ".tgz") {
			continue
		}

		releasefh, err := f.Open()
		if err != nil {
			return errors.Wrap(err, "opening embedded file")
		}

		err = a.analyzeReleaseTarball(records, f.Name, releasefh)
		if err != nil {
			releasefh.Close()

			return errors.Wrapf(err, "analyzing %s", f.Name)
		}

		releasefh.Close()

	}

	return nil
}

func (a *analysisGenerator) analyzeReleaseTarball(records analysis.Writer, artifact string, reader io.Reader) error {
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

		if path == "release.MF" {
			err = a.analyzeReleaseManifest(records, artifact, path, tarReader)
		} else if strings.HasPrefix(path, "jobs/") && strings.HasSuffix(path, ".tgz") {
			err = a.analyzeJobArtifactManifest(records, artifact, path, tarReader)
		} else {
			continue
		}

		if err != nil {
			return fmt.Errorf("analyzing artifact %s: %v", path, err)
		}
	}

	return nil
}

func (a *analysisGenerator) analyzeReleaseManifest(records analysis.Writer, releaseArtifact, artifact string, reader io.Reader) error {
	marshalBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.Wrap(err, "reading release.MF")
	}

	var spec result.RecordSpec

	err = yaml.Unmarshal(marshalBytes, &spec)
	if err != nil {
		return errors.Wrap(err, "parsing release.MF")
	}

	err = records.Write(result.Record{
		Release: releaseArtifact,
		Path:    artifact,
		Raw:     string(marshalBytes),
		Parsed:  safejson(spec).(result.RecordSpec),
	})
	if err != nil {
		return errors.Wrap(err, "writing result")
	}

	return nil
}

func (a *analysisGenerator) analyzeJobArtifactManifest(records analysis.Writer, releaseArtifact, artifact string, reader io.Reader) error {
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

		var spec result.RecordSpec

		err = yaml.Unmarshal(marshalBytes, &spec)
		if err != nil {
			return errors.Wrap(err, "parsing job.MF")
		}

		err = records.Write(result.Record{
			Release: releaseArtifact,
			Path:    artifact,
			Raw:     string(marshalBytes),
			Parsed:  safejson(spec).(result.RecordSpec),
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
