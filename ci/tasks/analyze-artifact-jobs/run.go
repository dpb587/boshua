package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Analysis struct {
	Artifact string      `json:"artifact"`
	Path     string      `json:"path"`
	Result   interface{} `json:"result"`
}

func main() {
	paths, err := filepath.Glob(os.Args[1])
	if err != nil {
		log.Fatalf("globbing: %v", err)
	} else if len(paths) != 1 {
		log.Fatalf("found %d matches for %s", len(paths), os.Args[1])
	}

	fh, err := os.Open(paths[0])
	if err != nil {
		log.Fatalf("opening file: %v", err)
	}

	defer fh.Close()

	gzReader, err := gzip.NewReader(fh)
	if err != nil {
		log.Fatalf("starting gzip: %v", err)
	}

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("advancing tar: %v", err)
		} else if header.Typeflag == tar.TypeDir {
			continue
		}

		path := strings.TrimPrefix(header.Name, "./")

		if !strings.HasSuffix(path, ".tgz") {
			continue
		} else if !strings.HasPrefix(path, "jobs/") {
			continue
		}

		analyzeArtifact(path, tarReader)
	}
}

func analyzeArtifact(artifact string, reader io.Reader) {
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		log.Fatalf("starting gzip: %v", err)
	}

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("advancing tar: %v", err)
		} else if header.Typeflag == tar.TypeDir {
			continue
		}

		path := strings.TrimPrefix(header.Name, "./")

		if path != "job.MF" {
			continue
		}

		specBytes, err := ioutil.ReadAll(tarReader)
		if err != nil {
			log.Fatalf("reading job spec: %v", err)
		}

		var spec map[interface{}]interface{}

		err = yaml.Unmarshal(specBytes, &spec)
		if err != nil {
			log.Fatalf("unmarshaling yaml: %v", err)
		}

		writeResult(artifact, path, safejson(spec))
	}
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

func writeResult(artifact, path string, analysis interface{}) {
	bytes, err := json.Marshal(Analysis{
		Artifact: artifact,
		Path:     path,
		Result:   analysis,
	})
	if err != nil {
		log.Fatalf("marshaling result: %v", err)
	}

	fmt.Printf("%s\n", bytes)
}
