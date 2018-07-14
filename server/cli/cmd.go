package cli

import (
	"net/http"

	"github.com/dpb587/boshua/cli/cmd/opts"
	releaseversionv2 "github.com/dpb587/boshua/releaseversion/api/v2/server"
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

	scheduler, err := c.AppOpts.GetScheduler()
	if err != nil {
		return errors.Wrap(err, "loading scheduler")
	}

	r := mux.NewRouter()

	r.HandleFunc("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("pong")) })).Methods(http.MethodGet)

	{
		datastore, err := c.AppOpts.GetReleaseIndex("default")
		if err != nil {
			return errors.Wrap(err, "loading release index")
		}

		releaseversionv2.NewHandler(c.AppOpts.GetLogger(), datastore, scheduler).RegisterHandlers(r)
	}

	return http.ListenAndServe(cfg.Bind, r)
}
