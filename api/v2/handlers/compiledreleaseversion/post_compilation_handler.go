package compiledreleaseversion

import (
	"fmt"
	"net/http"
	"reflect"

	api "github.com/dpb587/boshua/api/v2/models/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion/manager"
	"github.com/dpb587/boshua/compiledreleaseversion/task/compilation"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/scheduler"
	"github.com/dpb587/boshua/scheduler/concourse"
	"github.com/sirupsen/logrus"
)

type POSTCompilationHandler struct {
	logger                        logrus.FieldLogger
	cc                            *concourse.Runner
	compiledReleaseVersionManager *manager.Manager
	compiledReleaseVersionIndex   datastore.Index
}

func NewPOSTCompilationHandler(
	logger logrus.FieldLogger,
	cc *concourse.Runner,
	compiledReleaseVersionManager *manager.Manager,
	compiledReleaseVersionIndex datastore.Index,
) http.Handler {
	return &POSTCompilationHandler{
		logger: logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(POSTCompilationHandler{}).PkgPath(),
			"api.version":   "v2",
			"api.handler":   "compiledreleaseversion/compilation",
		}),
		cc: cc,
		compiledReleaseVersionManager: compiledReleaseVersionManager,
		compiledReleaseVersionIndex:   compiledReleaseVersionIndex,
	}
}

func (h *POSTCompilationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var status scheduler.Status

	baseLogger := applyLoggerContext(h.logger, r)

	releaseVersionRef, osVersionRef, logger, err := parseRequest(baseLogger, r)
	if err != nil {
		writeFailure(baseLogger, w, r, http.StatusBadRequest, fmt.Errorf("parsing request: %v", err))

		return
	}

	reqDataRef := compiledreleaseversion.Reference{
		ReleaseVersion: releaseVersionRef,
		OSVersion:      osVersionRef,
	}

	_, err = h.compiledReleaseVersionIndex.Find(reqDataRef)
	if err == datastore.MissingErr {
		releaseVersion, osVersion, err := h.compiledReleaseVersionManager.Resolve(reqDataRef)
		if err != nil {
			status := http.StatusInternalServerError

			if err == releaseversiondatastore.MissingErr || err == osversiondatastore.MissingErr {
				status = http.StatusBadRequest
			}

			writeFailure(logger, w, r, status, fmt.Errorf("resolving reference: %v", err))

			return
		}

		task := compilation.New(releaseVersion, osVersion)

		// check existing status
		status, err = h.cc.Status(task)
		if err != nil {
			writeFailure(logger, w, r, http.StatusInternalServerError, fmt.Errorf("checking task status: %v", err))

			return
		} else if status == scheduler.StatusUnknown {
			err = h.cc.Schedule(task)
			if err != nil {
				writeFailure(logger, w, r, http.StatusInternalServerError, fmt.Errorf("scheduling task: %v", err))

				return
			}

			// TODO log about scheduling

			status = scheduler.StatusPending
		}
	} else if err != nil {
		writeFailure(logger, w, r, http.StatusInternalServerError, fmt.Errorf("finding compiled release: %v", err))

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
			// TODO handle other errors?
			complete = true
		}
	case scheduler.StatusFailed:
		complete = true
	}

	writeResponse(logger, w, r, api.POSTCompilationResponse{
		Status:   string(status),
		Complete: complete,
	})
}
