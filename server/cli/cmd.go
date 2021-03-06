package cli

import (
	"net/http"

	"github.com/dpb587/boshua/server/handlers"
	// releaseversionv2 "github.com/dpb587/boshua/releaseversion/api/v2/server"
	"github.com/dpb587/boshua/config"
	"github.com/dpb587/boshua/config/provider/setter"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Cmd struct {
	setter.AppConfig `no-flag:"true"`
}

func (c *Cmd) Execute(extra []string) error {
	c.AppConfig.AppendLoggerFields(logrus.Fields{"cli.command": "server"})

	cfg := c.AppConfig.Config.Server

	r := mux.NewRouter()

	r.HandleFunc("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("pong")) })).Methods(http.MethodGet)

	if cfg.Mount.CLI != "" {
		r.PathPrefix("/cli/").Handler(http.StripPrefix("/cli/", http.FileServer(http.Dir(cfg.Mount.CLI)))).Methods(http.MethodGet)
	}

	if cfg.Mount.UI != "" {
		r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir(cfg.Mount.UI)))).Methods(http.MethodGet)
	}

	if cfg.Redirect.Root != "" {
		r.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, cfg.Redirect.Root, http.StatusFound)
		}))
	}

	releaseIndex, err := c.AppConfig.GetReleaseIndex(config.DefaultName)
	if err != nil {
		return errors.Wrap(err, "loading release index")
	}

	releaseComilationIndex, err := c.AppConfig.GetReleaseCompilationIndex(config.DefaultName)
	if err != nil {
		return errors.Wrap(err, "loading release index")
	}

	stemcellIndex, err := c.AppConfig.GetStemcellIndex(config.DefaultName)
	if err != nil {
		return errors.Wrap(err, "loading stemcell index")
	}

	scheduler, err := c.AppConfig.GetScheduler()
	if err != nil {
		return errors.Wrap(err, "loading scheduler")
	}

	handlers.NewGraphqlV2(
		c.AppConfig.GetLogger(),
		releaseIndex,
		c.AppConfig.GetReleaseAnalysisIndex,
		releaseComilationIndex,
		stemcellIndex,
		c.AppConfig.GetStemcellAnalysisIndex,
		scheduler,
	).Mount(r)

	downloader, err := c.AppConfig.GetDownloader()
	if err != nil {
		return errors.Wrap(err, "loading downloader")
	}

	handlers.NewMirrorsV2(
		c.AppConfig.GetLogger(),
		releaseIndex,
		c.AppConfig.GetReleaseAnalysisIndex,
		releaseComilationIndex,
		stemcellIndex,
		c.AppConfig.GetStemcellAnalysisIndex,
		downloader,
	).Mount(r)

	handlers.NewBOSHioV2(
		c.AppConfig.GetLogger(),
		releaseIndex,
		stemcellIndex,
	).Mount(r)

	return http.ListenAndServe(cfg.Bind, r)
}
