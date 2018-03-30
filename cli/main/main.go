package main

import (
	"net/http"
	"os"

	"github.com/dpb587/bosh-compiled-releases/api/v2/handlers"
	"github.com/dpb587/bosh-compiled-releases/compiler"
	compiledreleaseversionsaggregate "github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions/aggregate"
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions/legacybcr"
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions/presentbcr"
	releaseversionsaggregate "github.com/dpb587/bosh-compiled-releases/datastore/releaseversions/aggregate"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions/boshioreleaseindex"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions/boshmeta4releaseindex"
	stemcellversionsaggregate "github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions/aggregate"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions/boshiostemcellindex"
	"github.com/dpb587/bosh-compiled-releases/util"
	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	cc := compiler.Compiler{
		Target:       "dpb587-nightwatch-aws-use1",
		Insecure:     true,
		URL:          "https://concourse.nightwatch-aws-use1.dpb.io:4443",
		Team:         "main",
		Username:     "concourse",
		Password:     "0ac23mfhem569wpbau6r",
		PipelinePath: "/Users/dpb587/Projects/src/github.com/dpb587/bosh-compiled-releases/ci/compilation.yml",
		SecretsPath:  "/Users/dpb587/Projects/src/github.com/dpb587/bosh-compiled-releases/pipeline-vars.yml",
	}
	releaseIndex := releaseversionsaggregate.New(
		boshmeta4releaseindex.New("git+https://github.com/dpb587/openvpn-bosh-release.git//", "/Users/dpb587/Projects/src/github.com/dpb587/openvpn-bosh-release"),
		boshmeta4releaseindex.New("git+https://github.com/dpb587/ssoca-bosh-release.git//", "/Users/dpb587/Projects/src/github.com/dpb587/ssoca-bosh-release"),
		boshioreleaseindex.New("git+https://github.com/bosh-io/releases-index.git//", "/Users/dpb587/Projects/bosh-io/releases-index"),
	)
	stemcellIndex := stemcellversionsaggregate.New(
		boshiostemcellindex.New("git+https://github.com/bosh-io/stemcells-core-index.git//published/", "/Users/dpb587/Projects/bosh-io/stemcells-core-index/published"),
		boshiostemcellindex.New("git+https://github.com/bosh-io/stemcells-windows-index.git//published/", "/Users/dpb587/Projects/bosh-io/stemcells-windows-index/published"),
	)
	compiledReleaseIndex := compiledreleaseversionsaggregate.New(
		presentbcr.New(releaseIndex, "git@github.com:dpb587/bosh-compiled-releases-index.git", "/Users/dpb587/Projects/src/github.com/dpb587/bosh-compiled-releases-index"),
		legacybcr.New(releaseIndex, "git@github.com:dpb587/bosh-compiled-releases.git", "/Users/dpb587/Projects/src/github.com/dpb587/bosh-compiled-releases.gopath/src/github.com/dpb587/bosh-compiled-releases"),
	)
	releaseStemcellResolver := util.NewReleaseStemcellResolver(releaseIndex, stemcellIndex)

	r := mux.NewRouter()
	r.Handle("/v2/compiled-release-version/info", handlers.NewCRVInfoHandler(&cc, compiledReleaseIndex, releaseStemcellResolver)).Methods("POST")
	// r.Handle("/v2/compiled-release-version/log", handlers.NewCRVInfoHandler(compiledReleaseIndex)).Methods("POST")
	r.Handle("/v2/compiled-release-version/request", handlers.NewCRVRequestHandler(&cc, releaseStemcellResolver, compiledReleaseIndex)).Methods("POST")
	r.Handle("/v2/release-versions/list", handlers.NewRVListHandler(releaseIndex)).Methods("POST")
	// r.Handle("/v2/release-version/info", handlers.NewCRVInfoHandler(compiledReleaseIndex)).Methods("POST")
	// r.Handle("/v2/release-version/list-compiled-stemcells", handlers.NewCRVInfoHandler(compiledReleaseIndex)).Methods("POST")
	r.Handle("/v2/stemcell-versions/list", handlers.NewSVListHandler(stemcellIndex)).Methods("POST")
	// r.Handle("/v2/stemcell-version/info", handlers.NewCRVInfoHandler(compiledReleaseIndex)).Methods("POST")
	// r.Handle("/v2/stemcell-version/list-compiled-releases", handlers.NewCRVInfoHandler(compiledReleaseIndex)).Methods("POST")

	loggedRouter := gorillahandlers.LoggingHandler(os.Stdout, r)

	http.Handle("/", loggedRouter)
	http.ListenAndServe(":8080", nil)
}
