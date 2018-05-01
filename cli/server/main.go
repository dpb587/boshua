package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	analysisaggregate "github.com/dpb587/boshua/analysis/datastore/aggregate"
	analysisfactory "github.com/dpb587/boshua/analysis/datastore/factory"
	"github.com/dpb587/boshua/api/logging"
	handlersv2 "github.com/dpb587/boshua/api/v2/handlers"
	compiledreleaseversiondatastore "github.com/dpb587/boshua/compiledreleaseversion/datastore"
	compiledreleaseversionaggregate "github.com/dpb587/boshua/compiledreleaseversion/datastore/aggregate"
	compiledreleaseversionfactory "github.com/dpb587/boshua/compiledreleaseversion/datastore/factory"
	"github.com/dpb587/boshua/compiledreleaseversion/manager"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	osversionstemcellversionindex "github.com/dpb587/boshua/osversion/datastore/stemcellversionindex"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	releaseversionaggregate "github.com/dpb587/boshua/releaseversion/datastore/aggregate"
	releaseversionfactory "github.com/dpb587/boshua/releaseversion/datastore/factory"
	"github.com/dpb587/boshua/scheduler/concourse"
	"github.com/dpb587/boshua/server/config"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversionaggregate "github.com/dpb587/boshua/stemcellversion/datastore/aggregate"
	stemcellversionfactory "github.com/dpb587/boshua/stemcellversion/datastore/factory"
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

	cc := &concourse.Runner{
		Target:      serverConfig.Concourse.Target,
		Insecure:    serverConfig.Concourse.Insecure,
		URL:         serverConfig.Concourse.URL,
		Team:        serverConfig.Concourse.Team,
		Username:    serverConfig.Concourse.Username,
		Password:    serverConfig.Concourse.Password,
		SecretsPath: serverConfig.Concourse.SecretsPath,
	}

	var rv releaseversiondatastore.Index

	{
		var all []releaseversiondatastore.Index
		factory := releaseversionfactory.New(logger)

		for _, cfg := range serverConfig.Releases {
			idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options)
			if err != nil {
				log.Panicf("creating release version: %v", err)
			}

			all = append(all, idx)
		}

		rv = releaseversionaggregate.New(all...)
	}

	var sv stemcellversiondatastore.Index

	{
		var all []stemcellversiondatastore.Index
		factory := stemcellversionfactory.New(logger)

		for _, cfg := range serverConfig.Stemcells {
			idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options)
			if err != nil {
				log.Panicf("creating stemcell version: %v", err)
			}

			all = append(all, idx)
		}

		sv = stemcellversionaggregate.New(all...)
	}

	var ov osversiondatastore.Index = osversionstemcellversionindex.New(sv, logger)

	var crv compiledreleaseversiondatastore.Index

	{
		var all []compiledreleaseversiondatastore.Index
		factory := compiledreleaseversionfactory.New(logger, rv)

		for _, cfg := range serverConfig.CompiledReleases {
			idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options)
			if err != nil {
				log.Panicf("creating compiled release version: %v", err)
			}

			all = append(all, idx)
		}

		crv = compiledreleaseversionaggregate.New(all...)
	}

	var analysis analysisdatastore.Index

	{
		var all []analysisdatastore.Index
		factory := analysisfactory.New(logger, rv)

		for _, cfg := range serverConfig.Analysis {
			idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options)
			if err != nil {
				log.Panicf("creating analysis: %v", err)
			}

			all = append(all, idx)
		}

		analysis = analysisaggregate.New(all...)
	}

	r := mux.NewRouter()
	handlersv2.Mount(r.PathPrefix("/v2").Subrouter(), logger, cc, manager.NewManager(rv, ov), crv, rv, ov, sv, analysis)

	loggingRouter := logging.NewLogging(logger, r)
	loggerContextRouter := logging.NewLoggerContext(loggingRouter)

	http.Handle("/", loggerContextRouter)
	http.ListenAndServe(":8080", nil)
}
