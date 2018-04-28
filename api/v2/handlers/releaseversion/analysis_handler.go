package releaseversion

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/api/v2/handlers/analysisutil"
	"github.com/dpb587/boshua/api/v2/httputil"
	"github.com/dpb587/boshua/api/v2/urlutil"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/scheduler/concourse"
	"github.com/sirupsen/logrus"
)

const AnalysisHandlerInfoURI = "/release-version/analysis/info"
const AnalysisHandlerQueueURI = "/release-version/analysis/queue"

type pkg struct{}

func NewAnalysisHandler(
	logger logrus.FieldLogger,
	cc *concourse.Runner,
	analysisIndex datastore.Index,
	releaseVersionIndex releaseversiondatastore.Index,
) *analysisutil.AnalysisHandler {
	return analysisutil.NewAnalysisHandler(
		logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(pkg{}).PkgPath(),
			"api.version":   "v2",
			"api.handler":   "releaseversion/analysis",
		}),
		cc,
		analysisIndex,
		func(logger logrus.FieldLogger, r *http.Request) (analysis.Reference, logrus.FieldLogger, error) {
			releaseVersionRef, err := urlutil.ReleaseVersionRefFromParam(r)
			if err != nil {
				return analysis.Reference{}, nil, fmt.Errorf("parsing release version: %v", err)
			}

			analyzer, err := urlutil.AnalysisAnalyzerFromParam(r)
			if err != nil {
				return analysis.Reference{}, nil, fmt.Errorf("parsing analyzer: %v", err)
			}

			logger = logger.WithFields(logrus.Fields{
				"boshua.release.name":      releaseVersionRef.Name,
				"boshua.release.version":   releaseVersionRef.Version,
				"boshua.release.checksum":  releaseVersionRef.Checksums[0].String(),
				"boshua.analysis.analyzer": analyzer,
			})

			releaseVersion, err := releaseVersionIndex.Find(releaseVersionRef)
			if err != nil {
				httperr := httputil.NewError(err, http.StatusInternalServerError, "release version index failed")

				if err == releaseversiondatastore.MissingErr {
					httperr = httputil.NewError(err, http.StatusNotFound, "release version not found")
				}

				return analysis.Reference{}, logger, httperr
			}

			analysisRef := analysis.Reference{
				Artifact: releaseVersion,
				Analyzer: analyzer,
			}

			return analysisRef, logger, nil
		},
	)
}
