package handlers

import (
	"github.com/dpb587/boshua/api/v2/handlers/compiledreleaseversion"
	"github.com/dpb587/boshua/api/v2/handlers/releaseversions"
	"github.com/dpb587/boshua/api/v2/handlers/stemcellversions"
	compiledreleaseversionsds "github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion/manager"
	releaseversionsds "github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/scheduler/concourse"
	stemcellversionsds "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func Mount(
	router *mux.Router,
	logger logrus.FieldLogger,
	cc *concourse.Runner,
	compiledReleaseVersionManager *manager.Manager,
	compiledReleaseVersionIndex compiledreleaseversionsds.Index,
	releaseVersionIndex releaseversionsds.Index,
	stemcellVersionIndex stemcellversionsds.Index,
) {
	router.Handle("/compiled-release-version/info", compiledreleaseversion.NewInfoHandler(logger, compiledReleaseVersionIndex)).Methods("POST")
	// router.Handle("/compiled-release-version/log", compiledreleaseversion.NewCRVInfoHandler(compiledReleaseVersionIndex)).Methods("POST")
	router.Handle("/compiled-release-version/request", compiledreleaseversion.NewRequestHandler(logger, cc, compiledReleaseVersionManager, compiledReleaseVersionIndex)).Methods("POST")
	router.Handle("/release-versions/list", releaseversions.NewListHandler(logger, releaseVersionIndex)).Methods("POST")
	// router.Handle("/release-version/info", handlers.NewCRVInfoHandler(compiledReleaseVersionIndex)).Methods("POST")
	// router.Handle("/release-version/list-compiled-stemcells", handlers.NewCRVInfoHandler(compiledReleaseVersionIndex)).Methods("POST")
	router.Handle("/stemcell-versions/list", stemcellversions.NewListHandler(logger, stemcellVersionIndex)).Methods("POST")
	// router.Handle("/stemcell-version/info", handlers.NewCRVInfoHandler(compiledReleaseVersionIndex)).Methods("POST")
	// router.Handle("/stemcell-version/list-compiled-releases", handlers.NewCRVInfoHandler(compiledReleaseVersionIndex)).Methods("POST")
}
