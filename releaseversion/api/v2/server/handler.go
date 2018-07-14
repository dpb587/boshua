package server

import (
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/analysis"
	analysisv2 "github.com/dpb587/boshua/analysis/api/v2"
	analysisserver "github.com/dpb587/boshua/analysis/api/v2/server"
	"github.com/dpb587/boshua/api/v2/urlutil"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/api/v2"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/server/httputil"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type pkg struct{}

type Handler struct {
	Analysis *analysisserver.AnalysisHandler

	index  datastore.Index
	logger logrus.FieldLogger
}

func NewHandler(logger logrus.FieldLogger, index datastore.Index, taskScheduler scheduler.Scheduler) *Handler {
	logger = logger.WithFields(logrus.Fields{
		"build.package": reflect.TypeOf(pkg{}).PkgPath(),
		"api.version":   "v2",
	})

	h := &Handler{
		index:  index,
		logger: logger,
	}

	h.Analysis = h.newAnalysis(taskScheduler)

	return h
}

func (h *Handler) newAnalysis(taskScheduler scheduler.Scheduler) *analysisserver.AnalysisHandler {
	return analysisserver.NewAnalysisHandler(
		h.logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(pkg{}).PkgPath(),
			"api.version":   "v2",
			"api.handler":   "releaseversion/analysis",
		}),
		taskScheduler,
		h.index.GetAnalysisDatastore(),
		false,
		func(logger logrus.FieldLogger, r *http.Request) (analysis.Reference, logrus.FieldLogger, error) {
			subject, logger, err := h.parseRequest(r)
			if err != nil {
				return analysis.Reference{}, nil, httputil.NewError(err, http.StatusBadRequest, "parsing request")
			}

			analyzer, err := analysisv2.AnalysisAnalyzerFromParam(r)
			if err != nil {
				return analysis.Reference{}, nil, errors.Wrap(err, "parsing analyzer")
			}

			logger = logger.WithField("boshua.analysis.analyzer", analyzer)

			analysisRef := analysis.Reference{
				Subject:  subject,
				Analyzer: analyzer,
			}

			return analysisRef, logger, nil
		},
	)
}

func (h *Handler) parseRequest(r *http.Request) (releaseversion.Artifact, logrus.FieldLogger, error) {
	releaseVersionRef, err := urlutil.ReleaseVersionRefFromParam(r)
	if err != nil {
		return releaseversion.Artifact{}, h.logger, errors.Wrap(err, "parsing release version")
	}

	logger := h.logger.WithFields(logrus.Fields{
		"boshua.release.name":    releaseVersionRef.Name,
		"boshua.release.version": releaseVersionRef.Version,
	})

	if len(releaseVersionRef.Checksums) > 0 {
		logger = logger.WithField("boshua.release.checksum", releaseVersionRef.Checksums[0].String())
	}

	releaseVersions, err := h.index.Filter(releaseVersionRef)
	if err != nil {
		return releaseversion.Artifact{}, logger, httputil.NewError(err, http.StatusInternalServerError, "release version index failed")
	} else if len(releaseVersions) == 0 {
		return releaseversion.Artifact{}, logger, httputil.NewError(datastore.NoMatchErr, http.StatusNotFound, "release version not found")
	} else if len(releaseVersions) > 1 {
		return releaseversion.Artifact{}, logger, httputil.NewError(datastore.MultipleMatchErr, http.StatusBadRequest, "multiple release versions found")
	}

	return releaseVersions[0], logger, nil
}

func (h Handler) RegisterHandlers(r *mux.Router) {
	r.HandleFunc(v2.InfoPath, h.GetInfo).Methods(http.MethodGet)
	// r.HandleFunc(v2.AnalysisAnalyzersPath, h.Analysis.GetAnalyzers).Methods(http.MethodGet) # TODO not really analyzers of a single analysis
	r.HandleFunc(v2.AnalysisInfoPath, h.Analysis.GetInfo).Methods(http.MethodGet)
	r.HandleFunc(v2.AnalysisQueuePath, h.Analysis.PostQueue).Methods(http.MethodPost)
}
