package compiledreleaseversion

import (
	"fmt"
	"net/http"
	"reflect"

	api "github.com/dpb587/boshua/api/v2/models/compiledreleaseversion"
	schedulerapi "github.com/dpb587/boshua/api/v2/models/scheduler"
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

const CompilationHandlerInfoURI = "/compiled-release-version/compilation/info"
const CompilationHandlerQueueURI = "/compiled-release-version/compilation/queue"

type CompilationHandler struct {
	logger                        logrus.FieldLogger
	cc                            *concourse.Runner
	compiledReleaseVersionIndex   datastore.Index
	compiledReleaseVersionManager *manager.Manager
}

func NewCompilationHandler(
	logger logrus.FieldLogger,
	cc *concourse.Runner,
	compiledReleaseVersionIndex datastore.Index,
	compiledReleaseVersionManager *manager.Manager,
) *CompilationHandler {
	return &CompilationHandler{
		logger: logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(CompilationHandler{}).PkgPath(),
			"api.version":   "v2",
			"api.handler":   "compiledreleaseversion/info",
		}),
		cc: cc,
		compiledReleaseVersionIndex:   compiledReleaseVersionIndex,
		compiledReleaseVersionManager: compiledReleaseVersionManager,
	}
}

func (h *CompilationHandler) InfoGET(w http.ResponseWriter, r *http.Request) {
	baseLogger := applyLoggerContext(h.logger, r)

	releaseVersionRef, osVersionRef, logger, err := parseRequest(baseLogger, r)
	if err != nil {
		writeFailure(baseLogger, w, r, http.StatusBadRequest, fmt.Errorf("parsing request: %v", err))

		return
	}

	ref := compiledreleaseversion.Reference{
		ReleaseVersion: releaseVersionRef,
		OSVersion:      osVersionRef,
	}

	releaseVersion, osVersion, errResolve := h.compiledReleaseVersionManager.Resolve(ref)
	if errResolve != nil {
		status := http.StatusInternalServerError

		if errResolve == releaseversiondatastore.MissingErr || errResolve == osversiondatastore.MissingErr {
			status = http.StatusBadRequest
		}

		err = errResolve

		writeFailure(logger, w, r, status, fmt.Errorf("resolving ref: %v", err))

		return
	}

	ref = compiledreleaseversion.Reference{
		ReleaseVersion: releaseVersion.Reference,
		OSVersion:      osVersion.Reference,
	}

	result, err := h.compiledReleaseVersionIndex.Find(ref)
	if err != nil {
		status := http.StatusInternalServerError

		if err == datastore.MissingErr {
			// differentiate missing compilation vs invalid release/os
			status = http.StatusNotFound
		}

		writeFailure(logger, w, r, status, fmt.Errorf("finding compiled release: %v", err))

		return
	}

	logger.Infof("compiled release found")

	writeResponse(logger, w, r, api.GETCompilationResponse{
		Data: result.MetalinkFile,
	})
}

func (h *CompilationHandler) QueuePOST(w http.ResponseWriter, r *http.Request) {
	var status scheduler.Status

	baseLogger := applyLoggerContext(h.logger, r)

	releaseVersionRef, osVersionRef, logger, err := parseRequest(baseLogger, r)
	if err != nil {
		writeFailure(baseLogger, w, r, http.StatusBadRequest, fmt.Errorf("parsing request: %v", err))

		return
	}

	ref := compiledreleaseversion.Reference{
		ReleaseVersion: releaseVersionRef,
		OSVersion:      osVersionRef,
	}

	releaseVersion, osVersion, errResolve := h.compiledReleaseVersionManager.Resolve(ref)
	if errResolve != nil {
		status := http.StatusInternalServerError

		if errResolve == releaseversiondatastore.MissingErr || errResolve == osversiondatastore.MissingErr {
			status = http.StatusBadRequest
		}

		err = errResolve

		writeFailure(logger, w, r, status, fmt.Errorf("resolving ref: %v", err))

		return
	}

	ref = compiledreleaseversion.Reference{
		ReleaseVersion: releaseVersion.Reference,
		OSVersion:      osVersion.Reference,
	}

	_, err = h.compiledReleaseVersionIndex.Find(ref)
	if err == datastore.MissingErr {
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
		_, err = h.compiledReleaseVersionIndex.Find(ref)
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
		Data: schedulerapi.TaskStatus{
			Status:   string(status),
			Complete: complete,
		},
	})
}
