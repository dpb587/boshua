package compiledreleaseversion

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/api/v2/middleware"
	"github.com/dpb587/boshua/api/v2/models"
	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/sirupsen/logrus"
)

type InfoHandler struct {
	logger                      logrus.FieldLogger
	compiledReleaseVersionIndex datastore.Index
}

func NewInfoHandler(
	logger logrus.FieldLogger,
	compiledReleaseVersionIndex datastore.Index,
) http.Handler {
	return &InfoHandler{
		logger: logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(InfoHandler{}).PkgPath(),
			"api.version":   "v2",
			"api.handler":   "compiledreleaseversion/info",
		}),
		compiledReleaseVersionIndex: compiledReleaseVersionIndex,
	}
}

func (h *InfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := h.applyLoggerContext(r)

	reqData, err := h.readData(r)
	if err != nil {
		logger.WithField("error", err).Errorf("processing request body")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR: processing request body\n"))

		return
	}

	logger = logger.WithFields(logrus.Fields{
		"release.name":     reqData.ReleaseVersionRef.Name,
		"release.version":  reqData.ReleaseVersionRef.Version,
		"release.checksum": reqData.ReleaseVersionRef.Checksum,
		"stemcell.os":      reqData.StemcellVersionRef.OS,
		"stemcell.version": reqData.StemcellVersionRef.Version,
	})

	result, err := h.compiledReleaseVersionIndex.Find(compiledreleaseversion.Reference{
		ReleaseVersion: releaseversion.Reference{
			Name:      reqData.ReleaseVersionRef.Name,
			Version:   reqData.ReleaseVersionRef.Version,
			Checksums: checksum.ImmutableChecksums{reqData.ReleaseVersionRef.Checksum},
		},
		StemcellVersion: stemcellversion.Reference{
			OS:      reqData.StemcellVersionRef.OS,
			Version: reqData.StemcellVersionRef.Version,
		},
	})
	if err == datastore.MissingErr {
		logger.Infof("compiled release not found")

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("ERROR: compiled release version not found\n"))

		return
	} else if err != nil {
		logger.WithField("error", err).Errorf("finding compiled release version")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: finding compiled release version\n"))

		return
	}

	logger.Infof("compiled release found")

	res := models.CRVInfoResponse{
		Data: models.CRVInfoResponseData{
			ReleaseVersionRef:  reqData.ReleaseVersionRef,
			StemcellVersionRef: reqData.StemcellVersionRef,
			Artifact:           result.MetalinkFile,
		},
	}

	resBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		logger.WithField("error", err).Errorf("marshalling response")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: marshalling response\n"))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
	w.Write([]byte("\n"))
}

func (h *InfoHandler) readData(r *http.Request) (*models.CRVInfoRequestData, error) {
	var data models.CRVInfoRequest

	dataBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("reading: %v", err)
	}

	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling: %v", err)
	}

	return &data.Data, nil
}

func (h *InfoHandler) applyLoggerContext(r *http.Request) logrus.FieldLogger {
	logger := h.logger

	if context := r.Context().Value(middleware.LoggerContext); context != nil {
		logger = logger.WithFields(context.(logrus.Fields))
	}

	return logger
}
