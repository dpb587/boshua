package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions"
)

type SVListHandler struct {
	stemcellVersionIndex stemcellversions.Index
}

func NewSVListHandler(stemcellVersionIndex stemcellversions.Index) http.Handler {
	return &SVListHandler{
		stemcellVersionIndex: stemcellVersionIndex,
	}
}

func (h *SVListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	results, err := h.stemcellVersionIndex.List()
	if err != nil {
		log.Printf("listing stemcell versions: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR: listing stemcell versions\n"))

		return
	}

	res := models.SVListResponse{}

	for _, result := range results {
		res.Data = append(res.Data, models.StemcellRef{
			OS:      result.OS,
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
