package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	"github.com/dpb587/bosh-compiled-releases/compiler"
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
)

type CRVRequestHandler struct {
	cc                          *compiler.Compiler
	releaseVersionIndex         releaseversions.Index
	stemcellVersionIndex        stemcellversions.Index
	compiledReleaseVersionIndex compiledreleaseversions.Index
}

func NewCRVRequestHandler(
	cc *compiler.Compiler,
	releaseVersionIndex releaseversions.Index,
	stemcellVersionIndex stemcellversions.Index,
	compiledReleaseVersionIndex compiledreleaseversions.Index,
) http.Handler {
	return &CRVRequestHandler{
		cc:                          cc,
		releaseVersionIndex:         releaseVersionIndex,
		stemcellVersionIndex:        stemcellVersionIndex,
		compiledReleaseVersionIndex: compiledReleaseVersionIndex,
	}
}

func (h *CRVRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req models.CRVRequestRequest

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("reading request body: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: reading request body\n"))

		return
	}

	err = json.Unmarshal(reqBytes, &req)
	if err != nil {
		log.Printf("unmarshaling request body: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: unmarshaling request body\n"))

		return
	}

	_, err = h.compiledReleaseVersionIndex.Find(compiledreleaseversions.CompiledReleaseVersionRef{
		Release: releaseversions.ReleaseVersionRef{
			Name:    req.Data.Release.Name,
			Version: req.Data.Release.Version,
			Checksum: releaseversions.Checksum{
				Type:  req.Data.Release.Checksum.Type,
				Value: req.Data.Release.Checksum.Value,
			},
		},
		Stemcell: stemcellversions.StemcellVersionRef{
			OS:      req.Data.Stemcell.OS,
			Version: req.Data.Stemcell.Version,
		},
	})
	if err == compiledreleaseversions.MissingErr {
		release, err := h.releaseVersionIndex.Find(releaseversions.ReleaseVersionRef{
			Name:    req.Data.Release.Name,
			Version: req.Data.Release.Version,
			Checksum: releaseversions.Checksum{
				Type:  req.Data.Release.Checksum.Type,
				Value: req.Data.Release.Checksum.Value,
			},
		})
		if err != nil { // @todo MissingErr
			log.Printf("resolving release reference: %v", err)

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ERROR: resolving release reference\n"))

			return
		}

		stemcell, err := h.stemcellVersionIndex.Find(stemcellversions.StemcellVersionRef{
			OS:      req.Data.Stemcell.OS,
			Version: req.Data.Stemcell.Version,
		})
		if err != nil { // @todo MissingErr
			log.Printf("resolving stemcell reference: %v", err)

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ERROR: resolving stemcell reference\n"))

			return
		}

		err = h.cc.Schedule(release, stemcell)
		if err != nil {
			log.Printf("scheduling compiled release: %v", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ERROR: scheduling compiled release"))

			return
		}
	} else if err == nil {
		// already compiled; race condition
		// emulate pending
	} else {
		log.Printf("checking compiled release version: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: checking compiled release version\n"))

		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{
    "status": "pending"
}
`))
}
