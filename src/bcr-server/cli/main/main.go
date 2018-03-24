package main

import (
	"bcr-server/compiledreleaseversions/legacybcr"
	"bcr-server/releaseversions/boshioreleaseindex"
	"net/http"

	"bcr-server/api.v2/handlers"
	"github.com/gorilla/mux"
)

func main() {
	releaseIndex := boshioreleaseindex.New("git+https://github.com/bosh-io/releases-index.git//", "/Users/dpb587/Projects/bosh-io/releases-index")
	index := legacybcr.New(releaseIndex, "/Users/dpb587/Projects/dpb587/bosh-compiled-releases.gopath/src/github.com/dpb587/bosh-compiled-releases")

	r := mux.NewRouter()
	r.Handle("/v2/lookup", handlers.NewLookupHandler(index)).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
