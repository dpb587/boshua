package api

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/server/httputil"
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
	baseLogger := httputil.ApplyLoggerContext(h.logger, r)

	compiledReleaseVersionRef, logger, err := parseRequest(baseLogger, r)
	if err != nil {
		httputil.WriteFailure(baseLogger, w, r, httputil.NewError(errors.Wrap(err, "parsing request"), http.StatusBadRequest, "parsing request failed"))

		return
	}

	releaseVersion, osVersion, errResolve := h.compiledReleaseVersionManager.Resolve(compiledReleaseVersionRef)
	if errResolve != nil {
		httperr := httputil.NewError(errResolve, http.StatusInternalServerError, "resolving reference failed")

		switch errResolve {
		case osversiondatastore.MultipleMatchErr:
		case osversiondatastore.NoMatchErr:
		case releaseversiondatastore.MultipleMatchErr:
		case releaseversiondatastore.NoMatchErr:
			httperr = httputil.NewError(errResolve, http.StatusBadRequest, errResolve.Error())
		}

		httputil.WriteFailure(baseLogger, w, r, httperr)

		return
	}

	results, err := h.compiledReleaseVersionIndex.Filter(compiledreleaseversion.Reference{
		ReleaseVersion: releaseVersion.Reference,
		OSVersion:      osVersion.Reference,
	})
	if err != nil {
		httputil.WriteFailure(baseLogger, w, r, httputil.NewError(err, http.StatusInternalServerError, "compiled release index failed"))

		return
	} else if len(results) == 0 {
		httputil.WriteFailure(baseLogger, w, r, httputil.NewError(datastore.NoMatchErr, http.StatusNotFound, datastore.NoMatchErr.Error()))

		return
	} else if len(results) > 1 {
		httputil.WriteFailure(baseLogger, w, r, httputil.NewError(datastore.MultipleMatchErr, http.StatusBadRequest, datastore.MultipleMatchErr.Error()))

		return
	}

	result := results[0]

	logger.Infof("compiled release found")

	httputil.WriteResponse(logger, w, r, api.GETCompilationInfoResponse{
		Data: api.GETCompilationInfoResponseData{
			Artifact: result.MetalinkFile(),
		},
	})
}

func (h *CompilationHandler) QueuePOST(w http.ResponseWriter, r *http.Request) {
	var status scheduler.Status

	baseLogger := httputil.ApplyLoggerContext(h.logger, r)

	compiledReleaseVersionRef, logger, err := parseRequest(baseLogger, r)
	if err != nil {
		httputil.WriteFailure(baseLogger, w, r, httputil.NewError(errors.Wrap(err, "parsing request"), http.StatusBadRequest, "parsing request failed"))

		return
	}

	releaseVersion, osVersion, errResolve := h.compiledReleaseVersionManager.Resolve(compiledReleaseVersionRef)
	if errResolve != nil {
		httperr := httputil.NewError(err, http.StatusInternalServerError, "resolving reference failed")

		if errResolve == releaseversiondatastore.NoMatchErr {
			httperr = httputil.NewError(err, http.StatusBadRequest, "release version not found")
		} else if errResolve == osversiondatastore.NoMatchErr {
			httperr = httputil.NewError(err, http.StatusBadRequest, "os version not found")
		}

		httputil.WriteFailure(baseLogger, w, r, httperr)

		return
	}

	compiledReleaseVersionRef = compiledreleaseversion.Reference{
		ReleaseVersion: releaseVersion.Reference,
		OSVersion:      osVersion.Reference,
	}

	results, err := h.compiledReleaseVersionIndex.Filter(compiledReleaseVersionRef)
	if err != nil {
		httputil.WriteFailure(baseLogger, w, r, httputil.NewError(errors.Wrap(err, "finding compiled release"), http.StatusInternalServerError, "compiled release index failed"))

		return
	} else if len(results) == 0 {
		task := compilation.New(releaseVersion, osVersion)

		// check existing status
		status, err = h.cc.Status(task)
		if err != nil {
			httputil.WriteFailure(baseLogger, w, r, httputil.NewError(errors.Wrap(err, "checking task status"), http.StatusInternalServerError, "checking task status failed"))

			return
		} else if status == scheduler.StatusUnknown {
			err = h.cc.Schedule(task)
			if err != nil {
				httputil.WriteFailure(baseLogger, w, r, httputil.NewError(errors.Wrap(err, "scheduling task"), http.StatusInternalServerError, "scheduling task failed"))

				return
			}

			// TODO log about scheduling

			status = scheduler.StatusPending
		}
	} else {
		status = scheduler.StatusSucceeded
	}

	var complete bool

	switch status {
	case scheduler.StatusSucceeded:
		results, err = h.compiledReleaseVersionIndex.Filter(compiledReleaseVersionRef)
		if err != nil {
			httputil.WriteFailure(baseLogger, w, r, httputil.NewError(errors.Wrap(err, "filtering"), http.StatusInternalServerError, "compiled release index failed"))

			return
		}

		if len(results) == 0 {
			// propagation delay
			status = scheduler.StatusFinishing
		} else {
			// TODO handle other errors?
			complete = true
		}
	case scheduler.StatusFailed:
		complete = true
	}

	httputil.WriteResponse(logger, w, r, api.POSTCompilationQueueResponse{
		Data: schedulerapi.TaskStatus{
			Status:   string(status),
			Complete: complete,
		},
	})
}
