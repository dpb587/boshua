package legacybcr

import (
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions/inmemory"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type index struct {
	localPath string

	inmemory compiledreleaseversions.Index
}

func New(releaseVersionIndex releaseversions.Index, localPath string) compiledreleaseversions.Index {
	idx := &index{
		localPath: localPath,
	}

	idx.inmemory = inmemory.New(idx.loader, releaseVersionIndex)

	return idx
}

func (i *index) List() ([]compiledreleaseversions.CompiledReleaseVersion, error) {
	return i.inmemory.List()
}

func (i *index) Find(ref compiledreleaseversions.CompiledReleaseVersionRef) (compiledreleaseversions.CompiledReleaseVersion, error) {
	return i.inmemory.Find(ref)
}

func (i *index) loader() ([]compiledreleaseversions.CompiledReleaseVersion, error) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/data/**/**/**/bcr.json", i.localPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	var inmemory = []compiledreleaseversions.CompiledReleaseVersion{}

	for _, bcrPath := range paths {
		bcrBytes, err := ioutil.ReadFile(bcrPath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", bcrPath, err)
		}

		scanner := bufio.NewScanner(bytes.NewBuffer(bcrBytes))
		for scanner.Scan() {
			var record Record

			err = json.Unmarshal(scanner.Bytes(), &record)
			if err != nil {
				return nil, fmt.Errorf("unmarshalling %s: %v", bcrPath, err)
			}

			inmemory = append(inmemory, compiledreleaseversions.CompiledReleaseVersion{
				CompiledReleaseVersionRef: compiledreleaseversions.CompiledReleaseVersionRef{
					Release: releaseversions.ReleaseVersionRef{
						Name:    record.Name,
						Version: record.Version,
						Checksum: releaseversions.Checksum{
							Type:  "sha1",
							Value: record.Source.Digest,
						},
					},
					Stemcell: stemcellversions.StemcellVersionRef{
						OS:      record.Stemcell.OS,
						Version: record.Stemcell.Version,
					},
				},
				TarballChecksums: releaseversions.Checksums{
					releaseversions.Checksum{
						Type:  "sha1",
						Value: record.Tarball.Digest,
					},
				},
				TarballURL: record.Tarball.URL,
			})
		}
	}

	return inmemory, nil
}
