package handlers

import (
	"encoding/json"
	"fmt"
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

func NewCRVInfoHandler(
	compiledReleaseVersionIndex compiledreleaseversions.Index,
) http.Handler {
	return &CRVInfoHandler{
		compiledReleaseVersionIndex: compiledReleaseVersionIndex,
	}
}

func (h *CRVInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqData, err := h.readData(r)
	if err != nil {
		log.Printf("processing request body: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: processing request body\n"))

		return
	}

	result, err := h.compiledReleaseVersionIndex.Find(compiledreleaseversions.CompiledReleaseVersionRef{
		Release: releaseversions.ReleaseVersionRef{
			Name:     reqData.Release.Name,
			Version:  reqData.Release.Version,
			Checksum: releaseversions.Checksum(reqData.Release.Checksum),
		},
		Stemcell: stemcellversions.StemcellVersionRef{
			OS:      reqData.Stemcell.OS,
			Version: reqData.Stemcell.Version,
		},
	})
	if err == compiledreleaseversions.MissingErr {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("ERROR: compiled release version not found\n"))

		return
	} else if err != nil {
		log.Printf("finding compiled release version: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: finding compiled release version\n"))

		return
	}

	var checksums []models.Checksum

	for _, checksum := range result.TarballChecksums {
		if checksum.Algorithm() != "sha1" && checksum.Algorithm() != "sha256" {
			continue
		}

		checksums = append(checksums, models.Checksum(checksum))
	}

	res := models.CRVInfoResponse{
		Data: models.CRVInfoResponseData{
			Release:  reqData.Release,
			Stemcell: reqData.Stemcell,
			Tarball: models.CRVInfoResponseDataCompiled{
				URL:       result.TarballURL,
				Size:      result.TarballSize,
				Published: result.TarballPublished,
				Checksums: checksums,
			},
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

func (h *CRVInfoHandler) readData(r *http.Request) (*models.CRVInfoRequestData, error) {
	var data models.CRVInfoRequest

	dataBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("reading: %v", err)
	}

	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling: %v", err)
	}

	return &data.Data, nil
}
