package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	handlersv2 "github.com/dpb587/boshua/api/v2/handlers"
	"github.com/dpb587/boshua/api/v2/middleware"
	"github.com/dpb587/boshua/compiler"
	"github.com/dpb587/boshua/datastore/compiledreleaseversions"
	compiledreleaseversionsaggregate "github.com/dpb587/boshua/datastore/compiledreleaseversions/aggregate"
	compiledreleaseversionsfactory "github.com/dpb587/boshua/datastore/compiledreleaseversions/factory"
	"github.com/dpb587/boshua/datastore/releaseversions"
	releaseversionsaggregate "github.com/dpb587/boshua/datastore/releaseversions/aggregate"
	releaseversionsfactory "github.com/dpb587/boshua/datastore/releaseversions/factory"
	"github.com/dpb587/boshua/datastore/stemcellversions"
	stemcellversionsaggregate "github.com/dpb587/boshua/datastore/stemcellversions/aggregate"
	stemcellversionsfactory "github.com/dpb587/boshua/datastore/stemcellversions/factory"
	"github.com/dpb587/boshua/server/config"
	"github.com/dpb587/boshua/util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	var logger = logrus.New()
	logger.Out = os.Stdout
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.DebugLevel

	serverConfigBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Panicf("reading config: %v", err)
	}

	var serverConfig config.Config
	err = yaml.Unmarshal(serverConfigBytes, &serverConfig)
	if err != nil {
		log.Panicf("parsing config: %v", err)
	}

	cc := &compiler.Compiler{
		Target:       serverConfig.Concourse.Target,
		Insecure:     serverConfig.Concourse.Insecure,
		URL:          serverConfig.Concourse.URL,
		Team:         serverConfig.Concourse.Team,
		Username:     serverConfig.Concourse.Username,
		Password:     serverConfig.Concourse.Password,
		PipelinePath: serverConfig.Concourse.PipelinePath,
		SecretsPath:  serverConfig.Concourse.SecretsPath,
	}

	var rv releaseversions.Index

	{
		var all []releaseversions.Index
		factory := releaseversionsfactory.New(logger)

		for _, cfg := range serverConfig.Releases {
			idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options)
			if err != nil {
				log.Panicf("creating release version: %v", err)
			}

			all = append(all, idx)
		}

		rv = releaseversionsaggregate.New(all...)
	}

	var sv stemcellversions.Index

	{
		var all []stemcellversions.Index
		factory := stemcellversionsfactory.New(logger)

		for _, cfg := range serverConfig.Stemcells {
			idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options)
			if err != nil {
				log.Panicf("creating stemcell version: %v", err)
			}

			all = append(all, idx)
		}

		sv = stemcellversionsaggregate.New(all...)
	}

	var crv compiledreleaseversions.Index

	{
		var all []compiledreleaseversions.Index
		factory := compiledreleaseversionsfactory.New(logger, rv)

		for _, cfg := range serverConfig.CompiledReleases {
			idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options)
			if err != nil {
				log.Panicf("creating compiled release version: %v", err)
			}

			all = append(all, idx)
		}

		crv = compiledreleaseversionsaggregate.New(all...)
	}

	r := mux.NewRouter()
	handlersv2.Mount(r.PathPrefix("/v2").Subrouter(), logger, cc, util.NewReleaseStemcellResolver(rv, sv), crv, rv, sv)

	loggingRouter := middleware.NewLogging(logger, r)
	loggerContextRouter := middleware.NewLoggerContext(loggingRouter)

	http.Handle("/", loggerContextRouter)
	http.ListenAndServe(":8080", nil)
}
