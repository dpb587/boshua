package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/dpb587/bosh-compiled-releases/api/v2/middleware"
	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
	"github.com/sirupsen/logrus"
)

type CRVInfoHandler struct {
	logger                      logrus.FieldLogger
	compiledReleaseVersionIndex compiledreleaseversions.Index
}

func NewCRVInfoHandler(
	logger logrus.FieldLogger,
	compiledReleaseVersionIndex compiledreleaseversions.Index,
) http.Handler {
	return &CRVInfoHandler{
		logger: logger.WithFields(logrus.Fields{
			"package":     reflect.TypeOf(CRVInfoHandler{}).PkgPath(),
			"api.version": "v2",
			"api.handler": "crv_info",
		}),
		compiledReleaseVersionIndex: compiledReleaseVersionIndex,
	}
}

func (h *CRVInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := h.applyLoggerContext(r)

	reqData, err := h.readData(r)
	if err != nil {
		logger.WithField("error", err).Errorf("processing request body")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: processing request body\n"))

		return
	}

	logger = logger.WithFields(logrus.Fields{
		"release.name":     reqData.Release.Name,
		"release.version":  reqData.Release.Version,
		"release.checksum": reqData.Release.Checksum,
		"stemcell.os":      reqData.Stemcell.OS,
		"stemcell.version": reqData.Stemcell.Version,
	})

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
		logger.Infof("compiled release not found")

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("ERROR: compiled release version not found\n"))

		return
	} else if err != nil {
		logger.WithField("error", err).Errorf("finding compiled release version")

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

	logger.Infof("compiled release found")

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
		logger.WithField("error", err).Errorf("marshalling response")

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

func (h *CRVInfoHandler) applyLoggerContext(r *http.Request) logrus.FieldLogger {
	logger := h.logger

	if context := r.Context().Value(middleware.LoggerContext); context != nil {
		logger = logger.WithFields(context.(logrus.Fields))
	}

	return logger
}
