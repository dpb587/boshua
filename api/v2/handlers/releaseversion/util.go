package releaseversion

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dpb587/boshua/api/v2/middleware"
	"github.com/dpb587/boshua/api/v2/urlutil"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/sirupsen/logrus"
)

func parseRequest(logger logrus.FieldLogger, r *http.Request) (releaseversion.Reference, string, logrus.FieldLogger, error) {
	releaseVersionRef, err := urlutil.ReleaseVersionRefFromParam(r)
	if err != nil {
		return releaseversion.Reference{}, "", nil, fmt.Errorf("parsing release version: %v", err)
	}

	paramValue, ok := r.URL.Query()["analysis.analyzer"]
	if !ok {
		return releaseversion.Reference{}, "", nil, fmt.Errorf("parameter 'analysis': missing")
	} else if len(paramValue) != 1 {
		return releaseversion.Reference{}, "", nil, fmt.Errorf("parameter 'analysis': %v", fmt.Errorf("expected 1 value, but found %d", len(paramValue)))
	} else if len(paramValue[0]) == 0 {
		return releaseversion.Reference{}, "", nil, fmt.Errorf("parameter 'analysis': %v", errors.New("expected non-empty value"))
	}

	analyzer := paramValue[0]

	logger = logger.WithFields(logrus.Fields{
		"boshua.release.name":      releaseVersionRef.Name,
		"boshua.release.version":   releaseVersionRef.Version,
		"boshua.release.checksum":  releaseVersionRef.Checksums[0].String(),
		"boshua.analysis.analyzer": analyzer,
	})

	return releaseVersionRef, analyzer, logger, nil
}

func writeFailure(logger logrus.FieldLogger, w http.ResponseWriter, r *http.Request, status int, err error) {
	logger.WithField("error", err).Errorf("request failed")

	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf("ERROR: %v\n", err))) // TODO obfuscate potentially sensitive errors
}

func writeResponse(logger logrus.FieldLogger, w http.ResponseWriter, r *http.Request, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		writeFailure(logger, w, r, http.StatusInternalServerError, fmt.Errorf("marshaling response: %v", err))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(dataBytes)
	w.Write([]byte("\n"))
}

func applyLoggerContext(logger logrus.FieldLogger, r *http.Request) logrus.FieldLogger {
	if context := r.Context().Value(middleware.LoggerContext); context != nil {
		logger = logger.WithFields(context.(logrus.Fields))
	}

	return logger
}
