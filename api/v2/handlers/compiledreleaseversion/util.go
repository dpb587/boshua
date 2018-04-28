package compiledreleaseversion

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dpb587/boshua/api/logging"
	"github.com/dpb587/boshua/api/v2/urlutil"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/sirupsen/logrus"
)

func parseRequest(logger logrus.FieldLogger, r *http.Request) (releaseversion.Reference, osversion.Reference, logrus.FieldLogger, error) {
	releaseVersionRef, err := urlutil.ReleaseVersionRefFromParam(r)
	if err != nil {
		return releaseversion.Reference{}, osversion.Reference{}, nil, fmt.Errorf("parsing release version: %v", err)
	}

	osVersionRef, err := urlutil.OSVersionRefFromParam(r)
	if err != nil {
		return releaseversion.Reference{}, osversion.Reference{}, nil, fmt.Errorf("parsing os version: %v", err)
	}

	logger = logger.WithFields(logrus.Fields{
		"boshua.release.name":     releaseVersionRef.Name,
		"boshua.release.version":  releaseVersionRef.Version,
		"boshua.release.checksum": releaseVersionRef.Checksums[0].String(),
		"boshua.os.name":          osVersionRef.Name,
		"boshua.os.version":       osVersionRef.Version,
	})

	return releaseVersionRef, osVersionRef, logger, nil
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
	if context := r.Context().Value(logging.LoggerContext); context != nil {
		logger = logger.WithFields(context.(logrus.Fields))
	}

	return logger
}
