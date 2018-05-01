package handlers

import (
	"net/http"

	analysisds "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/api/v2/handlers/compiledreleaseversion"
	"github.com/dpb587/boshua/api/v2/handlers/osversions"
	"github.com/dpb587/boshua/api/v2/handlers/releaseversion"
	"github.com/dpb587/boshua/api/v2/handlers/releaseversions"
	"github.com/dpb587/boshua/api/v2/handlers/stemcellversion"
	compiledreleaseversionds "github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion/manager"
	osversionds "github.com/dpb587/boshua/osversion/datastore"
	releaseversionds "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/scheduler/concourse"
	stemcellversionds "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func Mount(
	router *mux.Router,
	logger logrus.FieldLogger,
	cc *concourse.Runner,
	compiledReleaseVersionManager *manager.Manager,
	compiledReleaseVersionIndex compiledreleaseversionds.Index,
	releaseVersionIndex releaseversionds.Index,
	osVersionIndex osversionds.Index,
	stemcellVersionIndex stemcellversionds.Index,
	analysisIndex analysisds.Index,
) {
	{
		handler := releaseversion.NewAnalysisHandler(logger, cc, analysisIndex, releaseVersionIndex)

		router.HandleFunc(releaseversion.AnalysisHandlerInfoURI, handler.InfoGET).Methods(http.MethodGet)
		router.HandleFunc(releaseversion.AnalysisHandlerQueueURI, handler.QueuePOST).Methods(http.MethodPost)
	}

	{
		handler := stemcellversion.NewAnalysisHandler(logger, cc, analysisIndex, stemcellVersionIndex)

		router.HandleFunc(stemcellversion.AnalysisHandlerInfoURI, handler.InfoGET).Methods(http.MethodGet)
		router.HandleFunc(stemcellversion.AnalysisHandlerQueueURI, handler.QueuePOST).Methods(http.MethodPost)
	}

	{
		handler := compiledreleaseversion.NewAnalysisHandler(logger, cc, analysisIndex, compiledReleaseVersionIndex)

		router.HandleFunc(compiledreleaseversion.AnalysisHandlerInfoURI, handler.InfoGET).Methods(http.MethodGet)
		router.HandleFunc(compiledreleaseversion.AnalysisHandlerQueueURI, handler.QueuePOST).Methods(http.MethodPost)
	}

	router.Handle("/compiled-release-version/compilation", compiledreleaseversion.NewGETCompilationHandler(logger, compiledReleaseVersionManager, compiledReleaseVersionIndex)).Methods(http.MethodGet)
	// router.Handle("/compiled-release-version/log", compiledreleaseversion.NewCRVInfoHandler(compiledReleaseVersionIndex)).Methods(http.MethodPost)
	router.Handle("/compiled-release-version/compilation", compiledreleaseversion.NewPOSTCompilationHandler(logger, cc, compiledReleaseVersionManager, compiledReleaseVersionIndex)).Methods(http.MethodPost)
	router.Handle("/release-versions/list", releaseversions.NewListHandler(logger, releaseVersionIndex)).Methods(http.MethodPost)
	// router.Handle("/release-version/info", handlers.NewCRVInfoHandler(compiledReleaseVersionIndex)).Methods(http.MethodPost)
	// router.Handle("/release-version/list-compiled-stemcells", handlers.NewCRVInfoHandler(compiledReleaseVersionIndex)).Methods(http.MethodPost)
	router.Handle("/stemcell-versions/list", osversions.NewListHandler(logger, osVersionIndex)).Methods(http.MethodPost)
	// router.Handle("/stemcell-version/info", handlers.NewCRVInfoHandler(compiledReleaseVersionIndex)).Methods(http.MethodPost)
	// router.Handle("/stemcell-version/list-compiled-releases", handlers.NewCRVInfoHandler(compiledReleaseVersionIndex)).Methods(http.MethodPost)
}
