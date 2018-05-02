package analysisutil

import (
	"fmt"
	"net/http"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/analysis/task"
	"github.com/dpb587/boshua/api/v2/httputil"
	api "github.com/dpb587/boshua/api/v2/models/analysis"
	schedulerapi "github.com/dpb587/boshua/api/v2/models/scheduler"
	"github.com/dpb587/boshua/scheduler"
	"github.com/dpb587/boshua/scheduler/concourse"
	"github.com/sirupsen/logrus"
)

type AnalysisHandlerRequestParser func(logger logrus.FieldLogger, r *http.Request) (analysis.Reference, logrus.FieldLogger, error)

type AnalysisHandler struct {
	logger          logrus.FieldLogger
	cc              *concourse.Runner
	analysisIndex   datastore.Index
	privilegedTasks bool
	requestParser   AnalysisHandlerRequestParser
}

func NewAnalysisHandler(
	logger logrus.FieldLogger,
	cc *concourse.Runner,
	analysisIndex datastore.Index,
	privilegedTasks bool,
	requestParser AnalysisHandlerRequestParser,
) *AnalysisHandler {
	return &AnalysisHandler{
		logger:          logger,
		cc:              cc,
		analysisIndex:   analysisIndex,
		privilegedTasks: privilegedTasks,
		requestParser:   requestParser,
	}
}

func (h *AnalysisHandler) InfoGET(w http.ResponseWriter, r *http.Request) {
	baseLogger := httputil.ApplyLoggerContext(h.logger, r)

	analysisRef, logger, err := h.requestParser(baseLogger, r)
	if err != nil {
		httputil.WriteFailure(baseLogger, w, r, httputil.NewError(err, http.StatusBadRequest, "request parsing failed"))

		return
	} else if err = h.validateAnalyzer(analysisRef); err != nil {
		httputil.WriteFailure(baseLogger, w, r, err)

		return
	}

	result, err := h.analysisIndex.Find(analysisRef)
	if err != nil {
		httperr := httputil.NewError(err, http.StatusInternalServerError, "analysis index failed")

		if err == datastore.MissingErr {
			httperr = httputil.NewError(err, http.StatusNotFound, "analysis not found")
		}

		httputil.WriteFailure(logger, w, r, httperr)

		return
	}

	logger.Infof("analysis found")

	httputil.WriteResponse(logger, w, r, api.GETAnalysisResponse{
		Data: result.ArtifactMetalink().Files[0],
	})
}

func (h *AnalysisHandler) QueuePOST(w http.ResponseWriter, r *http.Request) {
	var status scheduler.Status

	baseLogger := httputil.ApplyLoggerContext(h.logger, r)

	analysisRef, logger, err := h.requestParser(baseLogger, r)
	if err != nil {
		httputil.WriteFailure(baseLogger, w, r, httputil.NewError(err, http.StatusBadRequest, "request parsing failed"))

		return
	} else if err = h.validateAnalyzer(analysisRef); err != nil {
		httputil.WriteFailure(baseLogger, w, r, err)

		return
	}

	_, err = h.analysisIndex.Find(analysisRef)
	if err == datastore.MissingErr {
		t := task.New(analysisRef.Artifact.(analysis.Subject), analysisRef.Analyzer, h.privilegedTasks)

		// check existing status
		status, err = h.cc.Status(t)
		if err != nil {
			httputil.WriteFailure(logger, w, r, httputil.NewError(fmt.Errorf("checking task status: %v", err), http.StatusInternalServerError, "checking task status failed"))

			return
		} else if status == scheduler.StatusUnknown {
			err = h.cc.Schedule(t)
			if err != nil {
				httputil.WriteFailure(logger, w, r, httputil.NewError(fmt.Errorf("scheduling task: %v", err), http.StatusInternalServerError, "scheduling task failed"))

				return
			}

			logger.Infof("analysis scheduled")

			status = scheduler.StatusPending
		}
	} else if err != nil {
		httputil.WriteFailure(logger, w, r, httputil.NewError(fmt.Errorf("finding analysis: %v", err), http.StatusInternalServerError, "analysis index failed"))

		return
	} else {
		status = scheduler.StatusSucceeded
	}

	var complete bool

	switch status {
	case scheduler.StatusSucceeded:
		_, err = h.analysisIndex.Find(analysisRef)
		if err == datastore.MissingErr {
			// haven't reloaded it yet; delay them
			status = scheduler.StatusFinishing
		} else {
			// TODO handle other errors?
			complete = true
		}
	case scheduler.StatusFailed:
		complete = true
	}

	httputil.WriteResponse(logger, w, r, api.POSTAnalysisResponse{
		Data: schedulerapi.TaskStatus{
			Status:   string(status),
			Complete: complete,
		},
	})
}

func (h *AnalysisHandler) validateAnalyzer(analysisRef analysis.Reference) error {
	subject, ok := analysisRef.Artifact.(analysis.Subject)
	if !ok {
		return fmt.Errorf("TODO panic about using bad subjects for analysis")
	}

	found := false

	for _, supportedAnalyzer := range subject.SupportedAnalyzers() {
		if supportedAnalyzer == analysisRef.Analyzer {
			found = true

			break
		}
	}

	if !found {
		err := fmt.Errorf("unsupported analyzer: %s", analysisRef.Analyzer)

		return httputil.NewError(err, http.StatusNotFound, err.Error())
	}

	return nil
}
