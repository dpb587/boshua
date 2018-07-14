package server

import (
	"fmt"
	"net/http"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	analysistask "github.com/dpb587/boshua/analysis/task"
	api "github.com/dpb587/boshua/api/v2/models/analysis"
	schedulerapi "github.com/dpb587/boshua/api/v2/models/scheduler"
	"github.com/dpb587/boshua/server/httputil"
	"github.com/dpb587/boshua/task"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type AnalysisHandlerRequestParser func(logger logrus.FieldLogger, r *http.Request) (analysis.Reference, logrus.FieldLogger, error)

type AnalysisHandler struct {
	logger          logrus.FieldLogger
	scheduler       scheduler.Scheduler
	analysisIndex   datastore.Index
	privilegedTasks bool
	requestParser   AnalysisHandlerRequestParser
}

func NewAnalysisHandler(
	logger logrus.FieldLogger,
	scheduler scheduler.Scheduler,
	analysisIndex datastore.Index,
	privilegedTasks bool,
	requestParser AnalysisHandlerRequestParser,
) *AnalysisHandler {
	return &AnalysisHandler{
		logger:          logger,
		scheduler:       scheduler,
		analysisIndex:   analysisIndex,
		privilegedTasks: privilegedTasks,
		requestParser:   requestParser,
	}
}

func (h *AnalysisHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	baseLogger := httputil.ApplyLoggerContext(h.logger, r)

	analysisRef, logger, err := h.requestParser(baseLogger, r)
	if err != nil {
		httputil.WriteFailure(baseLogger, w, r, httputil.NewError(err, http.StatusBadRequest, "request parsing failed"))

		return
	} else if err = h.validateAnalyzer(analysisRef); err != nil {
		httputil.WriteFailure(baseLogger, w, r, err)

		return
	}

	results, err := h.analysisIndex.Filter(analysisRef)
	if err != nil {
		httputil.WriteFailure(logger, w, r, httputil.NewError(err, http.StatusInternalServerError, "analysis index failed"))

		return
	} else if len(results) == 0 {
		httputil.WriteFailure(logger, w, r, httputil.NewError(datastore.NoMatchErr, http.StatusNotFound, datastore.NoMatchErr.Error()))

		return
	} else if len(results) > 1 {
		httputil.WriteFailure(logger, w, r, httputil.NewError(datastore.MultipleMatchErr, http.StatusBadRequest, datastore.MultipleMatchErr.Error()))

		return
	}

	logger.Infof("analysis found")

	httputil.WriteResponse(logger, w, r, api.GETInfoResponse{
		Data: api.GETInfoResponseData{
			Artifact: results[0].MetalinkFile(),
		},
	})
}

func (h *AnalysisHandler) PostQueue(w http.ResponseWriter, r *http.Request) {
	var status task.Status

	baseLogger := httputil.ApplyLoggerContext(h.logger, r)

	analysisRef, logger, err := h.requestParser(baseLogger, r)
	if err != nil {
		httputil.WriteFailure(baseLogger, w, r, httputil.NewError(err, http.StatusBadRequest, "request parsing failed"))

		return
	} else if err = h.validateAnalyzer(analysisRef); err != nil {
		httputil.WriteFailure(baseLogger, w, r, err)

		return
	}

	analyses, err := h.analysisIndex.Filter(analysisRef)
	if err == datastore.NoMatchErr {
		taskDefinition, err := analysistask.New(analysisRef.Subject.(analysis.Subject), analysisRef.Analyzer)
		if err != nil {
			httputil.WriteFailure(logger, w, r, httputil.NewError(errors.Wrap(err, "building task definition"), http.StatusInternalServerError, "building analysis task definition failed"))

			return
		}

		// check existing status
		t, err := h.scheduler.Schedule(taskDefinition)
		if err != nil {
			httputil.WriteFailure(logger, w, r, httputil.NewError(errors.Wrap(err, "scheduling task"), http.StatusInternalServerError, "scheduling task failed"))

			return
		}

		status, err := t.Status()
		if err != nil {
			httputil.WriteFailure(logger, w, r, httputil.NewError(errors.Wrap(err, "checking task status"), http.StatusInternalServerError, "scheduling task failed"))

			return
		} else if status == task.StatusUnknown {
			if err != nil {
				httputil.WriteFailure(logger, w, r, httputil.NewError(errors.Wrap(err, "scheduling task"), http.StatusInternalServerError, "scheduling task failed"))

				return
			}

			logger.Infof("analysis scheduled")

			status = task.StatusPending
		}
	} else if err != nil {
		httputil.WriteFailure(logger, w, r, httputil.NewError(errors.Wrap(err, "filtering"), http.StatusInternalServerError, "analysis index failed"))

		return
	} else if err != nil {
		httputil.WriteFailure(logger, w, r, httputil.NewError(errors.Wrap(err, "finding analysis"), http.StatusInternalServerError, "analysis index failed"))

		return
	} else {
		status = task.StatusSucceeded
	}

	var complete bool

	switch status {
	case task.StatusSucceeded:
		analyses, err = h.analysisIndex.Filter(analysisRef)
		if err != nil {
			httputil.WriteFailure(logger, w, r, httputil.NewError(errors.Wrap(err, "filtering"), http.StatusInternalServerError, "analysis index failed"))

			return
		} else if len(analyses) == 0 {
			// haven't reloaded it yet; delay them
			status = task.StatusFinishing
		} else {
			// TODO handle other errors?
			complete = true
		}
	case task.StatusFailed:
		complete = true
	}

	httputil.WriteResponse(logger, w, r, api.POSTQueueResponse{
		Data: schedulerapi.TaskStatus{
			Status:   string(status),
			Complete: complete,
		},
	})
}

func (h *AnalysisHandler) validateAnalyzer(analysisRef analysis.Reference) error {
	subject, ok := analysisRef.Subject.(analysis.Subject)
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
