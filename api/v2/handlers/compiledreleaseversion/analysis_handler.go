package compiledreleaseversion

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/api/v2/handlers/analysisutil"
	"github.com/dpb587/boshua/api/v2/httputil"
	"github.com/dpb587/boshua/api/v2/urlutil"
	compiledreleaseversiondatastore "github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/scheduler/concourse"
	"github.com/sirupsen/logrus"
)

const AnalysisHandlerInfoURI = "/compiled-release-version/analysis/info"
const AnalysisHandlerQueueURI = "/compiled-release-version/analysis/queue"

type pkg struct{}

func NewAnalysisHandler(
	logger logrus.FieldLogger,
	cc *concourse.Runner,
	analysisIndex datastore.Index,
	compiledReleaseVersionIndex compiledreleaseversiondatastore.Index,
) *analysisutil.AnalysisHandler {
	return analysisutil.NewAnalysisHandler(
		logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(pkg{}).PkgPath(),
			"api.version":   "v2",
			"api.handler":   "compiledreleaseversion/analysis",
		}),
		cc,
		analysisIndex,
		false,
		func(baseLogger logrus.FieldLogger, r *http.Request) (analysis.Reference, logrus.FieldLogger, error) {
			compiledReleaseVersionRef, logger, err := parseRequest(baseLogger, r)
			if err != nil {
				return analysis.Reference{}, nil, fmt.Errorf("parsing request: %v", err)
			}

			analyzer, err := urlutil.AnalysisAnalyzerFromParam(r)
			if err != nil {
				return analysis.Reference{}, nil, fmt.Errorf("parsing analyzer: %v", err)
			}

			logger = logger.WithField("boshua.analysis.analyzer", analyzer)

			compiledReleaseVersions, err := compiledReleaseVersionIndex.Filter(compiledReleaseVersionRef)
			if err != nil {
				return analysis.Reference{}, logger, httputil.NewError(err, http.StatusInternalServerError, "compiled release version index failed")
			} else if len(compiledReleaseVersions) == 0 {
				return analysis.Reference{}, logger, httputil.NewError(datastore.NoMatchErr, http.StatusNotFound, "compiled release version not found")
			} else if len(compiledReleaseVersions) > 1 {
				return analysis.Reference{}, logger, httputil.NewError(datastore.MultipleMatchErr, http.StatusBadRequest, "multiple compiled release versions found")
			}

			analysisRef := analysis.Reference{
				Artifact: compiledReleaseVersions[0],
				Analyzer: analyzer,
			}

			return analysisRef, logger, nil
		},
	)
}
