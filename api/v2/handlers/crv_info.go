package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
)

type CRVInfoHandler struct {
	compiledReleaseVersionIndex compiledreleaseversions.Index
}

func NewCRVInfoHandler(compiledReleaseVersionIndex compiledreleaseversions.Index) http.Handler {
	return &CRVInfoHandler{
		compiledReleaseVersionIndex: compiledReleaseVersionIndex,
	}
}

func (h *CRVInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req models.CRVInfoRequest

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

	result, err := h.compiledReleaseVersionIndex.Find(compiledreleaseversions.CompiledReleaseVersionRef{
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
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found\n"))

		return
	} else if err != nil {
		log.Printf("finding compiled release version: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: finding compiled release version\n"))

		return
	}

	var checksums []models.Checksum

	for _, checksum := range result.TarballChecksums {
		if checksum.Type != "sha1" && checksum.Type != "sha256" {
			continue
		}

		checksums = append(checksums, models.Checksum{
			Type:  checksum.Type,
			Value: checksum.Value,
		})
	}

	res := models.CRVInfoResponse{
		Data: models.CRVInfoResponseCompiledRelease{
			URL:       result.TarballURL,
			Checksums: checksums,
			Release:   req.Data.Release,
			Stemcell:  req.Data.Stemcell,
		},
	}

	resBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Printf("marshalling response: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: marshalling response\n"))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
	w.Write([]byte("\n"))
}
