package osversions

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/api/v2/models"
	"github.com/dpb587/boshua/osversion/datastore"
	"github.com/sirupsen/logrus"
)

type ListHandler struct {
	logger         logrus.FieldLogger
	osVersionIndex datastore.Index
}

func NewListHandler(logger logrus.FieldLogger, osVersionIndex datastore.Index) http.Handler {
	return &ListHandler{
		logger: logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(ListHandler{}).PkgPath(),
			"api.version":   "v2",
			"api.handler":   "osversions/list",
		}),
		osVersionIndex: osVersionIndex,
	}
}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	results, err := h.osVersionIndex.List()
	if err != nil {
		log.Printf("listing os versions: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: listing os versions\n"))

		return
	}

	res := models.OVListResponse{}

	for _, result := range results {
		res.Data = append(res.Data, models.OSVersionRef{
			Name:    result.Name,
			Version: result.Version,
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
