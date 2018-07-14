package server

import (
	"net/http"

	"github.com/dpb587/boshua/server/httputil"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/api/v2"
)

func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) {
	subject, logger, err := h.parseRequest(r)
	if err != nil {
		httputil.WriteFailure(h.logger, w, r, httputil.NewError(err, http.StatusBadRequest, "parsing request"))

		return
	}

	httputil.WriteResponse(logger, w, r, v2.InfoResponse{
		Data: v2.InfoResponseData{
			Reference: v2.FromReference(subject.Reference().(stemcellversion.Reference)),
			Artifact:  subject.MetalinkFile(),
		},
	})
}
