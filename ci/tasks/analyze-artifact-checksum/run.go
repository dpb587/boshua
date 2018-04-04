package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Analysis struct {
	Artifact string      `json:"artifact"`
	Path     string      `json:"path"`
	Result   interface{} `json:"result"`
}

type Checksums []Checksum

func (cs Checksums) Write(p []byte) (int, error) {
	for _, c := range cs {
		c.Write(p)
	}

	return len(p), nil // TODO optimistic
}

type Checksum struct {
	algorithm string
	hasher    hash.Hash
}

func (c Checksum) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%s:%x", c.algorithm, c.hasher.Sum(nil))), nil
}

func (c Checksum) Write(p []byte) (int, error) {
	return c.hasher.Write(p)
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

		checksums := newChecksums()

		_, err = io.Copy(checksums, tarReader)
		if err != nil {
			log.Fatalf("creating checksum: %v", err)
		}

		writeResult(artifact, path, checksums)
	}
}

func newChecksums() Checksums {
	return Checksums{
		Checksum{algorithm: "md5", hasher: md5.New()},
		Checksum{algorithm: "sha1", hasher: sha1.New()},
		Checksum{algorithm: "sha256", hasher: sha256.New()},
		Checksum{algorithm: "sha512", hasher: sha512.New()},
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
