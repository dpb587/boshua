package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Analysis struct {
	Artifact string      `json:"artifact"`
	Path     string      `json:"path"`
	Result   interface{} `json:"result"`
}

type File struct {
	Type       string     `json:"type"`
	Path       string     `json:"path"`
	Link       string     `json:"link,omitempty"`
	Size       int64      `json:"size,omitempty"`
	Mode       int64      `json:"mode"`
	Uid        int        `json:"uid"`
	Gid        int        `json:"gid"`
	Uname      string     `json:"uname"`
	Gname      string     `json:"gname"`
	ModTime    time.Time  `json:"modtime"`
	AccessTime *time.Time `json:"accesstime,omitempty"`
	ChangeTime *time.Time `json:"changetime,omitempty"`
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

		file := File{
			Type:    string(header.Typeflag),
			Path:    path,
			Link:    header.Linkname,
			Size:    header.Size,
			Mode:    header.Mode,
			Uid:     header.Uid,
			Gid:     header.Gid,
			Uname:   header.Uname,
			Gname:   header.Gname,
			ModTime: header.ModTime,
		}

		unknownTime := time.Time{}

		if header.Format == tar.FormatPAX || header.Format == tar.FormatGNU {
			if header.ChangeTime != unknownTime {
				file.ChangeTime = &header.ChangeTime
			}

			if header.AccessTime != unknownTime {
				file.AccessTime = &header.AccessTime
			}
		}

		writeResult(artifact, path, file)
	}
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

	log.Printf("%s\n", bytes)
}
