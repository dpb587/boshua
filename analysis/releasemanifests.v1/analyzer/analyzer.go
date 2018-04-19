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
	releasemanifests "github.com/dpb587/boshua/analysis/releasemanifests.v1"
	yaml "gopkg.in/yaml.v2"
)

type Analyzer struct {
	tarball string
}

var _ analysis.Analyzer = &Analyzer{}

func New(tarball string) Analyzer {
	return Analyzer{
		tarball: tarball,
	}
}

func (a Analyzer) Analyze(results analysis.Writer) error {
	fh, err := os.Open(a.tarball)
	if err != nil {
		return fmt.Errorf("opening file: %v", err)
	}

	defer fh.Close()

	gzReader, err := gzip.NewReader(fh)
	if err != nil {
		return fmt.Errorf("starting gzip: %v", err)
	}

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("advancing tar: %v", err)
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

func (a Analyzer) analyzeReleaseManifest(results analysis.Writer, artifact string, reader io.Reader) error {
	marshalBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("reading release.MF: %v", err)
	}

	var spec releasemanifests.RecordSpec

	err = yaml.Unmarshal(marshalBytes, &spec)
	if err != nil {
		return fmt.Errorf("parsing release.MF: %v", err)
	}

	err = results.Write(releasemanifests.Record{
		Path:     artifact,
		Manifest: safejson(spec).(releasemanifests.RecordSpec),
	})
	if err != nil {
		return fmt.Errorf("writing result: %v", err)
	}

	return nil
}

func (a Analyzer) analyzeJobArtifactManifest(results analysis.Writer, artifact string, reader io.Reader) error {
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("starting gzip: %v", err)
	}

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("advancing tar: %v", err)
		} else if header.Typeflag == tar.TypeDir {
			continue
		}

		path := strings.TrimPrefix(header.Name, "./")

		if path != "job.MF" {
			continue
		}

		marshalBytes, err := ioutil.ReadAll(tarReader)
		if err != nil {
			return fmt.Errorf("reading job.MF: %v", err)
		}

		var spec releasemanifests.RecordSpec

		err = yaml.Unmarshal(marshalBytes, &spec)
		if err != nil {
			return fmt.Errorf("parsing job.MF: %v", err)
		}

		err = results.Write(releasemanifests.Record{
			Path:     artifact,
			Manifest: safejson(spec).(releasemanifests.RecordSpec),
		})
		if err != nil {
			return fmt.Errorf("writing result: %v", err)
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
