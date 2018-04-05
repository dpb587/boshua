package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/dpb587/bosh-compiled-releases/api/v2/middleware"
	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	"github.com/dpb587/bosh-compiled-releases/compiler"
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
	"github.com/dpb587/bosh-compiled-releases/util"
	"github.com/sirupsen/logrus"
)

type CRVRequestHandler struct {
	logger                      logrus.FieldLogger
	cc                          *compiler.Compiler
	releaseStemcellResolver     *util.ReleaseStemcellResolver
	compiledReleaseVersionIndex compiledreleaseversions.Index
}

func NewCRVRequestHandler(
	logger logrus.FieldLogger,
	cc *compiler.Compiler,
	releaseStemcellResolver *util.ReleaseStemcellResolver,
	compiledReleaseVersionIndex compiledreleaseversions.Index,
) http.Handler {
	return &CRVRequestHandler{
		logger: logger.WithFields(logrus.Fields{
			"package":     reflect.TypeOf(CRVRequestHandler{}).PkgPath(),
			"api.version": "v2",
			"api.handler": "crv_request",
		}),
		cc: cc,
		releaseStemcellResolver:     releaseStemcellResolver,
		compiledReleaseVersionIndex: compiledReleaseVersionIndex,
	}
}

func (h *CRVRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	var status compiler.Status

	_, err = h.compiledReleaseVersionIndex.Find(compiledreleaseversions.CompiledReleaseVersionRef{
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
		release, stemcell, err := h.releaseStemcellResolver.Resolve(
			releaseversions.ReleaseVersionRef{
				Name:     reqData.Release.Name,
				Version:  reqData.Release.Version,
				Checksum: releaseversions.Checksum(reqData.Release.Checksum),
			},
			stemcellversions.StemcellVersionRef{
				OS:      reqData.Stemcell.OS,
				Version: reqData.Stemcell.Version,
			},
		)
		if err == releaseversions.MissingErr || err == stemcellversions.MissingErr {
			logger.WithField("error", err).Infof("resolving reference")

			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("not found: %s\n", err)))

			return
		} else if err != nil {
			logger.WithField("error", err).Errorf("resolving reference")

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ERROR: resolving reference\n"))

			return
		}

		// check existing status
		status, err = h.cc.Status(release, stemcell)
		if err != nil {
			logger.WithField("error", err).Errorf("checking compilation status")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ERROR: checking compilation status\n"))

			return
		} else if status == compiler.StatusUnknown {
			err = h.cc.Schedule(release, stemcell)
			if err != nil {
				logger.WithField("error", err).Errorf("scheduling compiled release")

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("ERROR: scheduling compiled release"))

				return
			}

			status = compiler.StatusPending
		}
	} else if err != nil {
		logger.WithField("error", err).Errorf("checking compiled release version")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: checking compiled release version\n"))

		return
	} else {
		status = compiler.StatusSucceeded
	}

	var complete bool

	switch status {
	case compiler.StatusSucceeded:
		_, err = h.compiledReleaseVersionIndex.Find(compiledreleaseversions.CompiledReleaseVersionRef{
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
			status = compiler.StatusFinishing
		} else {
			complete = true
		}
	case compiler.StatusFailed:
		complete = true
	}

	h.writeData(logger, w, r, models.CRVRequestResponse{
		Status:   string(status),
		Complete: complete,
	})
}

func (h *CRVRequestHandler) readData(r *http.Request) (*models.CRVRequestRequestData, error) {
	var data models.CRVRequestRequest

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

func (h *CRVRequestHandler) writeData(logger logrus.FieldLogger, w http.ResponseWriter, r *http.Request, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		logger.WithField("error", err).Errorf("processing response body")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: processing response body\n"))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(dataBytes)
	w.Write([]byte("\n"))
}

func (h *CRVRequestHandler) applyLoggerContext(r *http.Request) logrus.FieldLogger {
	logger := h.logger

	if context := r.Context().Value(middleware.LoggerContext); context != nil {
		logger = logger.WithFields(context.(logrus.Fields))
	}

	return logger
}
