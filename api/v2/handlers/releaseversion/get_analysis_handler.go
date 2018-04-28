package releaseversion

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	api "github.com/dpb587/boshua/api/v2/models/analysis"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/sirupsen/logrus"
)

type GETAnalysisHandler struct {
	logger              logrus.FieldLogger
	analysisIndex       datastore.Index
	releaseVersionIndex releaseversiondatastore.Index
}

func NewGETAnalysisHandler(
	logger logrus.FieldLogger,
	analysisIndex datastore.Index,
	releaseVersionIndex releaseversiondatastore.Index,
) http.Handler {
	return &GETAnalysisHandler{
		logger: logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(GETAnalysisHandler{}).PkgPath(),
			"api.version":   "v2",
			"api.handler":   "releaseversion/analysis",
		}),
		analysisIndex:       analysisIndex,
		releaseVersionIndex: releaseVersionIndex,
	}
}

func (h *GETAnalysisHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.analysisIndex.Find(analysis.Reference{
		Artifact: releaseVersion.Reference,
		Analyzer: analyzer,
	})
	if err != nil {
		status := http.StatusInternalServerError

		if err == datastore.MissingErr {
			status = http.StatusNotFound
		}

		writeFailure(logger, w, r, status, fmt.Errorf("finding release version: %v", err))

		return
	}

	logger.Infof("release analysis found")

	writeResponse(logger, w, r, api.GETAnalysisResponse{
		Data: result.MetalinkFile,
	})
}
