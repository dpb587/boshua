package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	handlersv2 "github.com/dpb587/boshua/api/v2/handlers"
	"github.com/dpb587/boshua/api/v2/middleware"
	compiledreleaseversiondatastore "github.com/dpb587/boshua/compiledreleaseversion/datastore"
	compiledreleaseversionsaggregate "github.com/dpb587/boshua/compiledreleaseversion/datastore/aggregate"
	compiledreleaseversionsfactory "github.com/dpb587/boshua/compiledreleaseversion/datastore/factory"
	"github.com/dpb587/boshua/compiledreleaseversion/manager"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	osversionsaggregate "github.com/dpb587/boshua/osversion/datastore/aggregate"
	osversionsfactory "github.com/dpb587/boshua/osversion/datastore/factory"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	releaseversionsaggregate "github.com/dpb587/boshua/releaseversion/datastore/aggregate"
	releaseversionsfactory "github.com/dpb587/boshua/releaseversion/datastore/factory"
	"github.com/dpb587/boshua/scheduler/concourse"
	"github.com/dpb587/boshua/server/config"
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

	var sv osversiondatastore.Index

	{
		var all []osversiondatastore.Index
		factory := osversionsfactory.New(logger)

		for _, cfg := range serverConfig.Stemcells {
			idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options)
			if err != nil {
				log.Panicf("creating stemcell version: %v", err)
			}

			all = append(all, idx)
		}

		sv = osversionsaggregate.New(all...)
	}

	var crv compiledreleaseversiondatastore.Index

	{
		var all []compiledreleaseversiondatastore.Index
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
	handlersv2.Mount(r.PathPrefix("/v2").Subrouter(), logger, cc, manager.NewManager(rv, sv), crv, rv, sv)

	loggingRouter := middleware.NewLogging(logger, r)
	loggerContextRouter := middleware.NewLoggerContext(loggingRouter)

	http.Handle("/", loggerContextRouter)
	http.ListenAndServe(":8080", nil)
}
