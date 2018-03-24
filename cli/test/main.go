package main

import (
	"fmt"
	"log"

	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions/legacybcr"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions/boshioreleaseindex"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
	stemcellaggregate "github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions/aggregate"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions/boshiostemcellindex"
	"github.com/dpb587/bosh-compiled-releases/scheduler"
)

func main() {
	mainRelease()
	// mainCompiledRelease()
	// mainStemcell()
	// mainPipeline()
}

func mainPipeline() {
	releaseIndex := boshioreleaseindex.New("git+https://github.com/bosh-io/releases-index.git//", "/Users/dpb587/Projects/bosh-io/releases-index")
	stemcellIndex := stemcellaggregate.New(
		boshiostemcellindex.New("git+https://github.com/bosh-io/stemcells-core-index.git//published/", "/Users/dpb587/Projects/bosh-io/stemcells-core-index/published"),
		boshiostemcellindex.New("git+https://github.com/bosh-io/stemcells-windows-index.git//published/", "/Users/dpb587/Projects/bosh-io/stemcells-windows-index/published"),
	)

	release, err := releaseIndex.Find(releaseversions.ReleaseVersionRef{
		Name:    "openvpn",
		Version: "4.2.1",
		Checksum: releaseversions.Checksum{
			Type:  "sha1",
			Value: "80ac03f2ba2e142e2e0cbe1b23f340ec88d91c39",
		},
	})
	if err != nil {
		log.Fatalf("releasing: %v", err)
	}

	stemcell, err := stemcellIndex.Find(stemcellversions.StemcellVersionRef{
		OS:      "ubuntu-trusty",
		Version: "3468.22",
	})
	if err != nil {
		log.Fatalf("stemcelling: %v", err)
	}

	scheduler.Plan(release, stemcell)
}

func mainCompiledRelease() {
	releaseIndex := boshioreleaseindex.New("git+https://github.com/bosh-io/releases-index.git//", "/Users/dpb587/Projects/bosh-io/releases-index")
	index := legacybcr.New(releaseIndex, "/Users/dpb587/Projects/dpb587/bosh-compiled-releases.gopath/src/github.com/dpb587/bosh-compiled-releases")

	// result, err := index.List()
	result, err := index.Find(compiledreleaseversions.CompiledReleaseVersionRef{
		Release: releaseversions.ReleaseVersionRef{
			Name:    "openvpn",
			Version: "4.2.0",
			Checksum: releaseversions.Checksum{
				Type:  "sha1",
				Value: "56db9bd30ab2aabf7cafdad516d79be939d5d739",
			},
		},
		Stemcell: stemcellversions.StemcellVersionRef{
			OS:      "ubuntu-trusty",
			Version: "3468.22",
		},
	})
	if err != nil {
		log.Fatalf("resulting: %v", err)
	}

	fmt.Printf("%#+v\n", result)
}

func mainStemcell() {
	index := stemcellaggregate.New(
		boshiostemcellindex.New("git+https://github.com/bosh-io/stemcells-core-index.git//published/", "/Users/dpb587/Projects/bosh-io/stemcells-core-index/published"),
		boshiostemcellindex.New("git+https://github.com/bosh-io/stemcells-windows-index.git//published/", "/Users/dpb587/Projects/bosh-io/stemcells-windows-index/published"),
	)
	// result, err := index.List()
	result, err := index.Find(stemcellversions.StemcellVersionRef{
		OS:      "windows-2016",
		Version: "1709.3",
	})
	if err != nil {
		log.Fatalf("resulting: %v", err)
	}

	fmt.Printf("%#+v\n", result)
}

func mainRelease() {
	index := boshioreleaseindex.New("git+https://github.com/bosh-io/releases-index.git//", "/Users/dpb587/Projects/bosh-io/releases-index")
	result, err := index.List()
	// result, err := index.Find(releaseversions.ReleaseVersionRef{
	// 	Name:    "openvpn",
	// 	Version: "4.2.0",
	// 	Checksum: releaseversions.Checksum{
	// 		Type:  "sha1",
	// 		Value: "8e8ca38d82acfe51714128ca61e7b2894db798de",
	// 	},
	// })
	if err != nil {
		log.Fatalf("resulting: %v", err)
	}

	fmt.Printf("%#+v\n", result)
}
