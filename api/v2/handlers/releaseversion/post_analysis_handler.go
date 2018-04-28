package releaseversion

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/analysis/task"
	api "github.com/dpb587/boshua/api/v2/models/analysis"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/scheduler"
	"github.com/dpb587/boshua/scheduler/concourse"
	"github.com/sirupsen/logrus"
)

type POSTAnalysisHandler struct {
	logger              logrus.FieldLogger
	cc                  *concourse.Runner
	analysisIndex       datastore.Index
	releaseVersionIndex releaseversiondatastore.Index
}

func NewPOSTAnalysisHandler(
	logger logrus.FieldLogger,
	cc *concourse.Runner,
	analysisIndex datastore.Index,
	releaseVersionIndex releaseversiondatastore.Index,
) http.Handler {
	return &POSTAnalysisHandler{
		logger: logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(POSTAnalysisHandler{}).PkgPath(),
			"api.version":   "v2",
			"api.handler":   "releaseversion/analysis",
		}),
		cc:                  cc,
		analysisIndex:       analysisIndex,
		releaseVersionIndex: releaseVersionIndex,
	}
}

func (h *POSTAnalysisHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var status scheduler.Status

	baseLogger := applyLoggerContext(h.logger, r)

	releaseVersionRef, analyzer, logger, err := parseRequest(baseLogger, r)
	if err != nil {
		writeFailure(baseLogger, w, r, http.StatusBadRequest, fmt.Errorf("parsing request: %v", err))

		return
	}

	releaseVersion, err := h.releaseVersionIndex.Find(releaseVersionRef)
	if err != nil {
		status := http.StatusInternalServerError

		if err == datastore.MissingErr {
			status = http.StatusNotFound
		}

		writeFailure(logger, w, r, status, fmt.Errorf("finding release version: %v", err))

		return
	}

	found := false

	for _, supportedAnalyzer := range releaseVersion.SupportedAnalyzers() {
		if supportedAnalyzer == analyzer {
			found = true

			break
		}
	}

	if !found {
		writeFailure(logger, w, r, http.StatusNotFound, fmt.Errorf("unsupported analyzer: %s", analyzer))

		return
	}

	analysisRef := analysis.Reference{
		Artifact: releaseVersion.Reference,
		Analyzer: analyzer,
	}

	_, err = h.analysisIndex.Find(analysisRef)
	if err == datastore.MissingErr {
		t := task.New(releaseVersion, analyzer)

		// check existing status
		status, err = h.cc.Status(t)
		if err != nil {
			writeFailure(logger, w, r, http.StatusInternalServerError, fmt.Errorf("checking task status: %v", err))

			return
		} else if status == scheduler.StatusUnknown {
			err = h.cc.Schedule(t)
			if err != nil {
				writeFailure(logger, w, r, http.StatusInternalServerError, fmt.Errorf("scheduling task: %v", err))

				return
			}

			// TODO log about scheduling

			status = scheduler.StatusPending
		}
	} else if err != nil {
		writeFailure(logger, w, r, http.StatusInternalServerError, fmt.Errorf("finding release: %v", err))

		return
	} else {
		status = scheduler.StatusSucceeded
	}

	var complete bool

	switch status {
	case scheduler.StatusSucceeded:
		_, err = h.analysisIndex.Find(analysisRef)
		if err == datastore.MissingErr {
			status = scheduler.StatusFinishing
		} else {
			// TODO handle other errors?
			complete = true
		}
	case scheduler.StatusFailed:
		complete = true
	}

	writeResponse(logger, w, r, api.POSTAnalysisResponse{
		Status:   string(status),
		Complete: complete,
	})
}
