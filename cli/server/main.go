package main

import (
	"net/http"
	"os"

	handlersv2 "github.com/dpb587/bosh-compiled-releases/api/v2/handlers"
	"github.com/dpb587/bosh-compiled-releases/api/v2/middleware"
	"github.com/dpb587/bosh-compiled-releases/compiler"
	compiledreleaseversionsaggregate "github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions/aggregate"
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions/legacybcr"
	"github.com/dpb587/bosh-compiled-releases/datastore/compiledreleaseversions/presentbcr"
	releaseversionsaggregate "github.com/dpb587/bosh-compiled-releases/datastore/releaseversions/aggregate"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions/boshioreleaseindex"
	"github.com/dpb587/bosh-compiled-releases/datastore/releaseversions/boshmeta4releaseindex"
	stemcellversionsaggregate "github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions/aggregate"
	"github.com/dpb587/bosh-compiled-releases/datastore/stemcellversions/boshiostemcellindex"
	"github.com/dpb587/bosh-compiled-releases/util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	var logger = logrus.New()
	logger.Out = os.Stdout
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.DebugLevel

	cc := &compiler.Compiler{
		Target:       "dpb587-nightwatch-aws-use1",
		Insecure:     true,
		URL:          "https://concourse.nightwatch-aws-use1.dpb.io:4443",
		Team:         "main",
		Username:     "concourse",
		Password:     "0ac23mfhem569wpbau6r",
		PipelinePath: "/Users/dpb587/Projects/src/github.com/dpb587/bosh-compiled-releases/ci/compilation.yml",
		SecretsPath:  "/Users/dpb587/Projects/src/github.com/dpb587/bosh-compiled-releases/pipeline-vars.yml",
	}
	releaseVersionIndex := releaseversionsaggregate.New(
		boshmeta4releaseindex.New(logger.WithField("datastore", "dpb587/openvpn-bosh-release"), "git+https://github.com/dpb587/openvpn-bosh-release.git//", "/Users/dpb587/Projects/src/github.com/dpb587/openvpn-bosh-release"),
		boshmeta4releaseindex.New(logger.WithField("datastore", "dpb587/ssoca-bosh-release"), "git+https://github.com/dpb587/ssoca-bosh-release.git//", "/Users/dpb587/Projects/src/github.com/dpb587/ssoca-bosh-release"),
		boshioreleaseindex.New(logger.WithField("datastore", "bosh-io/releases-index"), "git+https://github.com/bosh-io/releases-index.git//", "/Users/dpb587/Projects/src/github.com/bosh-io/releases-index"),
	)
	stemcellVersionIndex := stemcellversionsaggregate.New(
		boshiostemcellindex.New(logger.WithField("datastore", "bosh-io/stemcells-core-index"), "git+https://github.com/bosh-io/stemcells-core-index.git//published/", "/Users/dpb587/Projects/src/github.com/bosh-io/stemcells-core-index/published"),
		boshiostemcellindex.New(logger.WithField("datastore", "bosh-io/stemcells-windows-index"), "git+https://github.com/bosh-io/stemcells-windows-index.git//published/", "/Users/dpb587/Projects/src/github.com/bosh-io/stemcells-windows-index/published"),
	)
	compiledReleaseVersionIndex := compiledreleaseversionsaggregate.New(
		presentbcr.New(logger.WithField("datastore", "present"), releaseVersionIndex, "git@github.com:dpb587/bosh-compiled-releases-index.git", "/Users/dpb587/Projects/src/github.com/dpb587/bosh-compiled-releases-index"),
		legacybcr.New(logger.WithField("datastore", "legacy"), releaseVersionIndex, "git@github.com:dpb587/bosh-compiled-releases.git", "/Users/dpb587/Projects/src/github.com/dpb587/bosh-compiled-releases.gopath/src/github.com/dpb587/bosh-compiled-releases"),
	)
	releaseStemcellResolver := util.NewReleaseStemcellResolver(releaseVersionIndex, stemcellVersionIndex)

	r := mux.NewRouter()
	handlersv2.Mount(
		r.PathPrefix("/v2").Subrouter(),
		logger,
		cc,
		releaseStemcellResolver,
		compiledReleaseVersionIndex,
		releaseVersionIndex,
		stemcellVersionIndex,
	)

	loggingRouter := middleware.NewLogging(logger, r)
	loggerContextRouter := middleware.NewLoggerContext(loggingRouter)

	http.Handle("/", loggerContextRouter)
	http.ListenAndServe(":8080", nil)
}
