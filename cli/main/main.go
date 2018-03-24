package main

import (
	"net/http"

	"github.com/dpb587/bosh-compiled-releases/api/v2/handlers"
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions/legacybcr"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions/boshioreleaseindex"
	stemcellaggregate "github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions/aggregate"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions/boshiostemcellindex"
	"github.com/gorilla/mux"
)

func main() {
	releaseIndex := boshioreleaseindex.New("git+https://github.com/bosh-io/releases-index.git//", "/Users/dpb587/Projects/bosh-io/releases-index")
	stemcellIndex := stemcellaggregate.New(
		boshiostemcellindex.New("git+https://github.com/bosh-io/stemcells-core-index.git//published/", "/Users/dpb587/Projects/bosh-io/stemcells-core-index/published"),
		boshiostemcellindex.New("git+https://github.com/bosh-io/stemcells-windows-index.git//published/", "/Users/dpb587/Projects/bosh-io/stemcells-windows-index/published"),
	)
	compiledReleaseIndex := legacybcr.New(releaseIndex, "/Users/dpb587/Projects/dpb587/bosh-compiled-releases.gopath/src/github.com/dpb587/bosh-compiled-releases")

	r := mux.NewRouter()
	r.Handle("/v2/compiled-release-version/info", handlers.NewCRVInfoHandler(compiledReleaseIndex)).Methods("POST")
	// r.Handle("/v2/compiled-release-version/log", handlers.NewCRVInfoHandler(compiledReleaseIndex)).Methods("POST")
	// r.Handle("/v2/compiled-release-version/request", handlers.NewCRVInfoHandler(compiledReleaseIndex)).Methods("POST")
	r.Handle("/v2/release-versions/list", handlers.NewRVListHandler(releaseIndex)).Methods("POST")
	// r.Handle("/v2/release-version/info", handlers.NewCRVInfoHandler(compiledReleaseIndex)).Methods("POST")
	// r.Handle("/v2/release-version/list-compiled-stemcells", handlers.NewCRVInfoHandler(compiledReleaseIndex)).Methods("POST")
	r.Handle("/v2/stemcell-versions/list", handlers.NewSVListHandler(stemcellIndex)).Methods("POST")
	// r.Handle("/v2/stemcell-version/info", handlers.NewCRVInfoHandler(compiledReleaseIndex)).Methods("POST")
	// r.Handle("/v2/stemcell-version/list-compiled-releases", handlers.NewCRVInfoHandler(compiledReleaseIndex)).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
