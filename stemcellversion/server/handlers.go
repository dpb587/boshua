package server

import (
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/gorilla/mux"
)

type Handlers struct {
	index datastore.Index
}

func NewHandlers(index datastore.Index) *Handlers {
	return &Handlers{
		index: index,
	}
}

func (h *Handlers) Mount(m *mux.Router) {
	m.Handle("/api/v2/stemcell/download", NewDownloadHandler(h.index))
	// m.Handle("/api/v2/stemcell/analysis/download", NewAnalysisDownloadHandler(h.index, h.analysisIndex))
}
