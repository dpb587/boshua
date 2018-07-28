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

	cfg, err := c.AppOpts.GetServerConfig()
	if err != nil {
		return errors.Wrap(err, "getting server config")
	}

	r := mux.NewRouter()

	r.HandleFunc("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("pong")) })).Methods(http.MethodGet)
	r.PathPrefix("/webui/").Handler(http.StripPrefix("/webui/", http.FileServer(http.Dir("webui")))).Methods(http.MethodGet) // TODO path assumptions

	releaseIndex, err := c.AppOpts.GetReleaseIndex("default")
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

	stemcellversionserver.NewHandlers(stemcellIndex, analysisIndex).Mount(r)
	handlers.NewGraphqlV2(releaseIndex, stemcellIndex).Mount(r)

	return http.ListenAndServe(cfg.Bind, r)
}
