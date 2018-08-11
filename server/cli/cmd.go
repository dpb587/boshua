package cli

import (
	"net/http"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/cli/cmd/opts"
	"github.com/dpb587/boshua/server/handlers"
	// releaseversionv2 "github.com/dpb587/boshua/releaseversion/api/v2/server"
	stemcellversionserver "github.com/dpb587/boshua/stemcellversion/server"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Cmd struct {
	AppOpts *opts.Opts `no-flag:"true"`
}

func (c *Cmd) Execute(extra []string) error {
	c.AppOpts.ConfigureLogger("server")

	cfgProvider, err := c.AppOpts.GetConfig()
	if err != nil {
		return errors.Wrap(err, "getting server config")
	}

	cfg := cfgProvider.Config.Server

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

	releaseIndex, err := c.AppOpts.GetReleaseIndex("default")
	if err != nil {
		return errors.Wrap(err, "loading release index")
	}

	releaseComilationIndex, err := c.AppOpts.GetCompiledReleaseIndex("default")
	if err != nil {
		return errors.Wrap(err, "loading release index")
	}

	stemcellIndex, err := c.AppOpts.GetStemcellIndex("default")
	if err != nil {
		return errors.Wrap(err, "loading stemcell index")
	}

	analysisIndex, err := c.AppOpts.GetAnalysisIndex(analysis.Reference{}) // TODO
	if err != nil {
		return errors.Wrap(err, "loading analysis index")
	}

	scheduler, err := c.AppOpts.GetScheduler()
	if err != nil {
		return errors.Wrap(err, "loading scheduler")
	}

	stemcellversionserver.NewHandlers(stemcellIndex, analysisIndex).Mount(r)
	handlers.NewGraphqlV2(c.AppOpts.GetLogger(), releaseIndex, releaseComilationIndex, stemcellIndex, scheduler).Mount(r)

	return http.ListenAndServe(cfg.Bind, r)
}
