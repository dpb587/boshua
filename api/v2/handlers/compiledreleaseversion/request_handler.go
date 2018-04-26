package compiledreleaseversion

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/api/v2/middleware"
	"github.com/dpb587/boshua/api/v2/models"
	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion/manager"
	"github.com/dpb587/boshua/compiledreleaseversion/task/compilation"
	"github.com/dpb587/boshua/releaseversion"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/scheduler"
	"github.com/dpb587/boshua/scheduler/concourse"
	"github.com/dpb587/boshua/stemcellversion"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/sirupsen/logrus"
)

type RequestHandler struct {
	logger                        logrus.FieldLogger
	cc                            *concourse.Runner
	compiledReleaseVersionManager *manager.Manager
	compiledReleaseVersionIndex   datastore.Index
}

func NewRequestHandler(
	logger logrus.FieldLogger,
	cc *concourse.Runner,
	compiledReleaseVersionManager *manager.Manager,
	compiledReleaseVersionIndex datastore.Index,
) http.Handler {
	return &RequestHandler{
		logger: logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(RequestHandler{}).PkgPath(),
			"api.version":   "v2",
			"api.handler":   "compiledreleaseversion/request",
		}),
		cc: cc,
		compiledReleaseVersionManager: compiledReleaseVersionManager,
		compiledReleaseVersionIndex:   compiledReleaseVersionIndex,
	}
}

func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := h.applyLoggerContext(r)

	reqData, err := h.readData(r)
	if err != nil {
		logger.WithField("error", err).Errorf("processing request body")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: processing request body\n"))

		return
	}

	logger = logger.WithFields(logrus.Fields{
		"release.name":     reqData.ReleaseVersionRef.Name,
		"release.version":  reqData.ReleaseVersionRef.Version,
		"release.checksum": reqData.ReleaseVersionRef.Checksum,
		"stemcell.os":      reqData.StemcellVersionRef.OS,
		"stemcell.version": reqData.StemcellVersionRef.Version,
	})

	var status scheduler.Status

	reqDataRef := compiledreleaseversion.Reference{
		ReleaseVersion: releaseversion.Reference{
			Name:      reqData.ReleaseVersionRef.Name,
			Version:   reqData.ReleaseVersionRef.Version,
			Checksums: checksum.ImmutableChecksums{reqData.ReleaseVersionRef.Checksum},
		},
		StemcellVersion: stemcellversion.Reference{
			OS:      reqData.StemcellVersionRef.OS,
			Version: reqData.StemcellVersionRef.Version,
		},
	}

	_, err = h.compiledReleaseVersionIndex.Find(reqDataRef)
	if err == datastore.MissingErr {
		releaseVersion, stemcellVersion, err := h.compiledReleaseVersionManager.Resolve(reqDataRef)
		if err == releaseversiondatastore.MissingErr || err == stemcellversiondatastore.MissingErr {
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

		task := compilation.New(releaseVersion, stemcellVersion)

		// check existing status
		status, err = h.cc.Status(task)
		if err != nil {
			logger.WithField("error", err).Errorf("checking compilation status")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ERROR: checking compilation status\n"))

			return
		} else if status == scheduler.StatusUnknown {
			err = h.cc.Schedule(task)
			if err != nil {
				logger.WithField("error", err).Errorf("scheduling compiled release")

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("ERROR: scheduling compiled release"))

				return
			}

			status = scheduler.StatusPending
		}
	} else if err != nil {
		logger.WithField("error", err).Errorf("checking compiled release version")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: checking compiled release version\n"))

		return
	} else {
		status = scheduler.StatusSucceeded
	}

	var complete bool

	switch status {
	case scheduler.StatusSucceeded:
		_, err = h.compiledReleaseVersionIndex.Find(reqDataRef)
		if err == datastore.MissingErr {
			status = scheduler.StatusFinishing
		} else {
			complete = true
		}
	case scheduler.StatusFailed:
		complete = true
	}

	h.writeData(logger, w, r, models.CRVRequestResponse{
		Status:   string(status),
		Complete: complete,
	})
}

func (h *RequestHandler) readData(r *http.Request) (*models.CRVRequestRequestData, error) {
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

func (h *RequestHandler) writeData(logger logrus.FieldLogger, w http.ResponseWriter, r *http.Request, data interface{}) {
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

func (h *RequestHandler) applyLoggerContext(r *http.Request) logrus.FieldLogger {
	logger := h.logger

	if context := r.Context().Value(middleware.LoggerContext); context != nil {
		logger = logger.WithFields(context.(logrus.Fields))
	}

	return logger
}
