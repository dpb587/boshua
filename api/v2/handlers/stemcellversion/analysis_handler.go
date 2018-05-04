package stemcellversion

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/api/v2/handlers/analysisutil"
	"github.com/dpb587/boshua/api/v2/httputil"
	"github.com/dpb587/boshua/api/v2/urlutil"
	"github.com/dpb587/boshua/scheduler/concourse"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/sirupsen/logrus"
)

const AnalysisHandlerInfoURI = "/stemcell-version/analysis/info"
const AnalysisHandlerQueueURI = "/stemcell-version/analysis/queue"

type pkg struct{}

func NewAnalysisHandler(
	logger logrus.FieldLogger,
	cc *concourse.Runner,
	analysisIndex datastore.Index,
	stemcellVersionIndex stemcellversiondatastore.Index,
) *analysisutil.AnalysisHandler {
	return analysisutil.NewAnalysisHandler(
		logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(pkg{}).PkgPath(),
			"api.version":   "v2",
			"api.handler":   "stemcellversion/analysis",
		}),
		cc,
		analysisIndex,
		true,
		func(logger logrus.FieldLogger, r *http.Request) (analysis.Reference, logrus.FieldLogger, error) {
			subject, logger, err := parseRequest(logger, r, stemcellVersionIndex)
			if err != nil {
				return analysis.Reference{}, nil, httputil.NewError(err, http.StatusBadRequest, "parsing request")
			}

			analyzer, err := urlutil.AnalysisAnalyzerFromParam(r)
			if err != nil {
				return analysis.Reference{}, nil, fmt.Errorf("parsing analyzer: %v", err)
			}

			logger = logger.WithField("boshua.analysis.analyzer", analyzer)

			analysisRef := analysis.Reference{
				Artifact: subject,
				Analyzer: analyzer,
			}

			return analysisRef, logger, nil
		},
	)
}
