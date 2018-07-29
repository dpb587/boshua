package server

import (
	"net/http"

	"github.com/dpb587/boshua/artifact/server/servercommon"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/stemcellversion/server/params"
	"github.com/pkg/errors"
)

type DownloadHandler struct {
	servercommon.MirrorHandler

	index datastore.Index
}

var _ http.Handler = &DownloadHandler{}

func NewDownloadHandler(index datastore.Index) *DownloadHandler {
	return &DownloadHandler{
		index: index,
	}
}

func (h *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filterParams, err := params.FilterParamsFromQuery(r)
	if err != nil {
		panic(errors.Wrap(err, "parsing request")) // TODO !panic
	}

	results, err := h.index.GetArtifacts(filterParams)
	if err != nil {
		panic(errors.Wrap(err, "finding stemcell")) // TODO !panic
	}

	result, err := datastore.RequireSingleResult(results)
	if err != nil {
		panic(errors.Wrap(err, "finding stemcell")) // TODO !panic
	}

	h.MirrorHandler.ServeHTTPArtifact(w, r, result)
}
