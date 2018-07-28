package server

import (
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/gorilla/mux"
)

type Handlers struct {
	index         datastore.Index
	analysisIndex analysisdatastore.Index
}

func NewHandlers(index datastore.Index, analysisIndex analysisdatastore.Index) *Handlers {
	return &Handlers{
		index:         index,
		analysisIndex: analysisIndex,
	}
}

func (h *Handlers) Mount(m *mux.Router) {
	m.Handle("/api/v2/stemcell/download", NewDownloadHandler(h.index))
	// m.Handle("/api/v2/stemcell/analysis/download", NewAnalysisDownloadHandler(h.index, h.analysisIndex))
}
