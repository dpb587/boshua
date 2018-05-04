package stemcellversion

import (
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/api/v2/httputil"
	api "github.com/dpb587/boshua/api/v2/models/stemcellversion"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/sirupsen/logrus"
)

const InfoHandlerURI = "/stemcell-version/info"

type InfoHandler struct {
	logger               logrus.FieldLogger
	stemcellVersionIndex stemcellversiondatastore.Index
}

func NewInfoHandler(
	logger logrus.FieldLogger,
	stemcellVersionIndex stemcellversiondatastore.Index,
) *InfoHandler {
	return &InfoHandler{
		logger: logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(pkg{}).PkgPath(),
			"api.version":   "v2",
			"api.handler":   "stemcellversion/info",
		}),
		stemcellVersionIndex: stemcellVersionIndex,
	}
}

func (h *InfoHandler) GET(w http.ResponseWriter, r *http.Request) {
	subject, logger, err := parseRequest(h.logger, r, h.stemcellVersionIndex)
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
