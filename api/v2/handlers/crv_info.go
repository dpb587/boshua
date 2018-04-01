package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	"github.com/dpb587/bosh-compiled-releases/compiler"
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
	"github.com/dpb587/bosh-compiled-releases/util"
)

type CRVInfoHandler struct {
	cc                          *compiler.Compiler
	compiledReleaseVersionIndex compiledreleaseversions.Index
	releaseStemcellResolver     *util.ReleaseStemcellResolver
}

func NewCRVInfoHandler(
	cc *compiler.Compiler,
	compiledReleaseVersionIndex compiledreleaseversions.Index,
	releaseStemcellResolver *util.ReleaseStemcellResolver,
) http.Handler {
	return &CRVInfoHandler{
		cc: cc,
		releaseStemcellResolver:     releaseStemcellResolver,
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
			Name:     req.Data.Release.Name,
			Version:  req.Data.Release.Version,
			Checksum: releaseversions.Checksum(req.Data.Release.Checksum),
		},
		Stemcell: stemcellversions.StemcellVersionRef{
			OS:      req.Data.Stemcell.OS,
			Version: req.Data.Stemcell.Version,
		},
	})
	if err == compiledreleaseversions.MissingErr {
		release, stemcell, err := h.releaseStemcellResolver.Resolve(
			releaseversions.ReleaseVersionRef{
				Name:     req.Data.Release.Name,
				Version:  req.Data.Release.Version,
				Checksum: releaseversions.Checksum(req.Data.Release.Checksum),
			},
			stemcellversions.StemcellVersionRef{
				OS:      req.Data.Stemcell.OS,
				Version: req.Data.Stemcell.Version,
			},
		)
		if err == releaseversions.MissingErr || err == stemcellversions.MissingErr {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found\n"))

			return
		} else if err != nil {
			log.Printf("resolving references: %v", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ERROR: resolving references\n"))

			return
		}

		status, err := h.cc.Status(release, stemcell)
		if err != nil {
			log.Printf("checking compilation status: %v", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ERROR: checking compilation status\n"))

			return
		} else if status == compiler.StatusUnknown {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found\n"))

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"data":{"status":"%s"}}`, status)))

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
			Status:   models.CRVInfoStatusAvailable,
			Release:  req.Data.Release,
			Stemcell: req.Data.Stemcell,
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
