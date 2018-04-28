package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dpb587/boshua/api/logging"
	"github.com/sirupsen/logrus"
)

func WriteFailure(logger logrus.FieldLogger, w http.ResponseWriter, r *http.Request, err error) {
	logger.WithField("error", err).Errorf("request failed")

	var status = http.StatusInternalServerError
	var msg = ""

	if httperr, ok := err.(Error); ok {
		status = httperr.Status
		msg = httperr.PublicError
	}

	if msg == "" {
		msg = http.StatusText(status)
	}

	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf("ERROR: %v\n", msg)))
}

func WriteResponse(logger logrus.FieldLogger, w http.ResponseWriter, r *http.Request, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		WriteFailure(logger, w, r, fmt.Errorf("marshaling response: %v", err))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(dataBytes)
	w.Write([]byte("\n"))
}

func ApplyLoggerContext(logger logrus.FieldLogger, r *http.Request) logrus.FieldLogger {
	if context := r.Context().Value(logging.LoggerContext); context != nil {
		logger = logger.WithFields(context.(logrus.Fields))
	}

	return logger
}
