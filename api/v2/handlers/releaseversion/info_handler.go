package releaseversion

import (
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/api/v2/httputil"
	api "github.com/dpb587/boshua/api/v2/models/releaseversion"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/sirupsen/logrus"
)

const InfoHandlerURI = "/release-version/info"

type InfoHandler struct {
	logger              logrus.FieldLogger
	releaseVersionIndex releaseversiondatastore.Index
}

func NewInfoHandler(
	logger logrus.FieldLogger,
	releaseVersionIndex releaseversiondatastore.Index,
) *InfoHandler {
	return &InfoHandler{
		logger: logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(pkg{}).PkgPath(),
			"api.version":   "v2",
			"api.handler":   "releaseversion/analysis",
		}),
		releaseVersionIndex: releaseVersionIndex,
	}
}

func (h *InfoHandler) GET(w http.ResponseWriter, r *http.Request) {
	subject, logger, err := parseRequest(h.logger, r, h.releaseVersionIndex)
	if err != nil {
		httputil.WriteFailure(h.logger, w, r, httputil.NewError(err, http.StatusBadRequest, "parsing request"))

		return
	}

	httputil.WriteResponse(logger, w, r, api.InfoResponse{
		Data: api.InfoResponseData{
			Reference: api.FromReference(subject.Reference),
			Artifact:  subject.ArtifactMetalink().Files[0],
		},
	})
}
