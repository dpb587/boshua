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
	"github.com/dpb587/boshua/compiledreleaseversion"
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
		func(logger logrus.FieldLogger, r *http.Request) (analysis.Reference, logrus.FieldLogger, error) {
			releaseVersionRef, err := urlutil.ReleaseVersionRefFromParam(r)
			if err != nil {
				return analysis.Reference{}, nil, fmt.Errorf("parsing release version: %v", err)
			}

			osVersionRef, err := urlutil.OSVersionRefFromParam(r)
			if err != nil {
				return analysis.Reference{}, nil, fmt.Errorf("parsing os version: %v", err)
			}

			analyzer, err := urlutil.AnalysisAnalyzerFromParam(r)
			if err != nil {
				return analysis.Reference{}, nil, fmt.Errorf("parsing analyzer: %v", err)
			}

			logger = logger.WithFields(logrus.Fields{
				"boshua.release.name":      releaseVersionRef.Name,
				"boshua.release.version":   releaseVersionRef.Version,
				"boshua.release.checksum":  releaseVersionRef.Checksums[0].String(),
				"boshua.os.name":           osVersionRef.Name,
				"boshua.os.version":        osVersionRef.Version,
				"boshua.analysis.analyzer": analyzer,
			})

			compiledReleaseVersionRef := compiledreleaseversion.Reference{
				ReleaseVersion: releaseVersionRef,
				OSVersion:      osVersionRef,
			}

			releaseVersion, err := compiledReleaseVersionIndex.Find(compiledReleaseVersionRef)
			if err != nil {
				httperr := httputil.NewError(err, http.StatusInternalServerError, "compiled release version index failed")

				if err == compiledreleaseversiondatastore.MissingErr {
					httperr = httputil.NewError(err, http.StatusNotFound, "compiled release version not found")
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
