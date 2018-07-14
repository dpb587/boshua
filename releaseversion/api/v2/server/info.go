package server

import (
	"net/http"

	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/api/v2"
	"github.com/dpb587/boshua/server/httputil"
)

func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) {
	subject, logger, err := h.parseRequest(r)
	if err != nil {
		httputil.WriteFailure(h.logger, w, r, httputil.NewError(err, http.StatusBadRequest, "parsing request"))

		return
	}

	httputil.WriteResponse(logger, w, r, v2.InfoResponse{
		Data: v2.InfoResponseData{
			Reference: v2.FromReference(subject.Reference().(releaseversion.Reference)),
			Artifact:  subject.MetalinkFile(),
		},
	})
}
