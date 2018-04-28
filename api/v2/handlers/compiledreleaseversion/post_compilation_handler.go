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
