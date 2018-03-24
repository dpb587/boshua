package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"bcr-server/compiledreleaseversions"
	"bcr-server/releaseversions"
	"bcr-server/stemcellversions"

	"bcr-server/api.v2/models"
)

type LookupHandler struct {
	compiledReleaseIndex compiledreleaseversions.Index
}

func NewLookupHandler(compiledReleaseIndex compiledreleaseversions.Index) http.Handler {
	return &LookupHandler{
		compiledReleaseIndex: compiledReleaseIndex,
	}
}

func (h *LookupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req models.LookupRequest

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

	result, err := h.compiledReleaseIndex.Find(compiledreleaseversions.CompiledReleaseVersionRef{
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
		log.Printf("finding compiled release: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: finding compiled release\n"))

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

	resBytes, err := json.Marshal(models.LookupResponse{
		CompiledRelease: models.LookupResponseCompiledRelease{
			URL:       result.TarballURL,
			Checksums: checksums,
			Release:   req.Data.Release,
			Stemcell:  req.Data.Stemcell,
		},
	})

	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
}
