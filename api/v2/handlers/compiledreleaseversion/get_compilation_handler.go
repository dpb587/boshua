package compiledreleaseversion

import (
	"fmt"
	"net/http"
	"reflect"

	api "github.com/dpb587/boshua/api/v2/models/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion/manager"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/sirupsen/logrus"
)

type GETCompilationHandler struct {
	logger                        logrus.FieldLogger
	compiledReleaseVersionManager *manager.Manager
	compiledReleaseVersionIndex   datastore.Index
}

func NewGETCompilationHandler(
	logger logrus.FieldLogger,
	compiledReleaseVersionManager *manager.Manager,
	compiledReleaseVersionIndex datastore.Index,
) http.Handler {
	return &GETCompilationHandler{
		logger: logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(GETCompilationHandler{}).PkgPath(),
			"api.version":   "v2",
			"api.handler":   "compiledreleaseversion/info",
		}),
		compiledReleaseVersionManager: compiledReleaseVersionManager,
		compiledReleaseVersionIndex:   compiledReleaseVersionIndex,
	}
}

func (h *GETCompilationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
