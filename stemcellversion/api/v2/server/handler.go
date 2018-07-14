package server

import (
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/analysis"
	analysisv2 "github.com/dpb587/boshua/analysis/api/v2"
	analysisserver "github.com/dpb587/boshua/analysis/api/v2/server"
	"github.com/dpb587/boshua/api/v2/urlutil"
	"github.com/dpb587/boshua/server/httputil"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/api/v2"
	"github.com/dpb587/boshua/stemcellversion/datastore"
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
			"api.handler":   "stemcellversion/analysis",
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

func (h *Handler) parseRequest(r *http.Request) (stemcellversion.Artifact, logrus.FieldLogger, error) {
	stemcellVersionRef, err := urlutil.StemcellVersionRefFromParam(r)
	if err != nil {
		return stemcellversion.Artifact{}, h.logger, errors.Wrap(err, "parsing stemcell version")
	}

	logger := h.logger.WithFields(logrus.Fields{
		"boshua.stemcell.name":    stemcellVersionRef.Name,
		"boshua.stemcell.version": stemcellVersionRef.Version,
	})

	stemcellVersion, err := h.index.Find(stemcellVersionRef)
	if err != nil {
		httperr := httputil.NewError(err, http.StatusInternalServerError, "stemcell version index failed")

		if err == datastore.NoMatchErr {
			httperr = httputil.NewError(err, http.StatusNotFound, "stemcell version not found")
		}

		return stemcellversion.Artifact{}, logger, httperr
	}

	return stemcellVersion, logger, nil
}

func (h Handler) RegisterHandlers(r *mux.Router) {
	r.HandleFunc(v2.InfoPath, h.GetInfo).Methods(http.MethodGet)
	// r.HandleFunc(v2.AnalysisAnalyzersPath, h.Analysis.GetAnalyzers).Methods(http.MethodGet) # TODO not really analyzers of a single analysis
	r.HandleFunc(v2.AnalysisInfoPath, h.Analysis.GetInfo).Methods(http.MethodGet)
	r.HandleFunc(v2.AnalysisQueuePath, h.Analysis.PostQueue).Methods(http.MethodPost)
}
