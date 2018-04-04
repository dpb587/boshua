package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions"
	"github.com/sirupsen/logrus"
)

type RVListHandler struct {
	logger              logrus.FieldLogger
	releaseVersionIndex releaseversions.Index
}

func NewRVListHandler(logger logrus.FieldLogger, releaseVersionIndex releaseversions.Index) http.Handler {
	return &RVListHandler{
		logger:              logger.WithField("package", reflect.TypeOf(RVListHandler{}).PkgPath()),
		releaseVersionIndex: releaseVersionIndex,
	}
}

func (h *RVListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	results, err := h.releaseVersionIndex.List()
	if err != nil {
		log.Printf("listing release versions: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: listing release versions\n"))

		return
	}

	res := models.RVListResponse{}

	for _, result := range results {
		var checksums []models.Checksum

		for _, checksum := range result.Checksums {
			if checksum.Algorithm() != "sha1" && checksum.Algorithm() != "sha256" {
				continue
			}

			checksums = append(checksums, models.Checksum(checksum))
		}

		res.Data = append(res.Data, models.ReleaseRef{
			Name:     result.Name,
			Version:  result.Version,
			Checksum: checksums[0],
		})
	}

	resBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Printf("marshalling response: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: marshalling response\n"))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
	w.Write([]byte("\n"))
}
