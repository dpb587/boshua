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
			stemcellVersionRef, err := urlutil.StemcellVersionRefFromParam(r)
			if err != nil {
				return analysis.Reference{}, nil, fmt.Errorf("parsing stemcell version: %v", err)
			}

			analyzer, err := urlutil.AnalysisAnalyzerFromParam(r)
			if err != nil {
				return analysis.Reference{}, nil, fmt.Errorf("parsing analyzer: %v", err)
			}

			logger = logger.WithFields(logrus.Fields{
				"boshua.stemcell.iaas":       stemcellVersionRef.IaaS,
				"boshua.stemcell.hypervisor": stemcellVersionRef.Hypervisor,
				"boshua.stemcell.os":         stemcellVersionRef.OS,
				"boshua.stemcell.version":    stemcellVersionRef.Version,
				"boshua.analysis.analyzer":   analyzer,
			})

			stemcellVersion, err := stemcellVersionIndex.Find(stemcellVersionRef)
			if err != nil {
				httperr := httputil.NewError(err, http.StatusInternalServerError, "stemcell version index failed")

				if err == stemcellversiondatastore.MissingErr {
					httperr = httputil.NewError(err, http.StatusNotFound, "stemcell version not found")
				}

				return analysis.Reference{}, logger, httperr
			}

			analysisRef := analysis.Reference{
				Artifact: stemcellVersion,
				Analyzer: analyzer,
			}

			return analysisRef, logger, nil
		},
	)
}
