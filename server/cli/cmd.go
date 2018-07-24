package cli

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dpb587/boshua/cli/cmd/opts"
	releaseversiongraphql "github.com/dpb587/boshua/releaseversion/graphql"
	stemcellversiongraphql "github.com/dpb587/boshua/stemcellversion/graphql"
	// releaseversionv2 "github.com/dpb587/boshua/releaseversion/api/v2/server"
	// stemcellversionv2 "github.com/dpb587/boshua/stemcellversion/api/v2/server"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
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

	// scheduler, err := c.AppOpts.GetScheduler()
	// if err != nil {
	// 	return errors.Wrap(err, "loading scheduler")
	// }

	r := mux.NewRouter()

	r.HandleFunc("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("pong")) })).Methods(http.MethodGet)
	r.PathPrefix("/webui/").Handler(http.StripPrefix("/webui/", http.FileServer(http.Dir("webui")))).Methods(http.MethodGet) // TODO path assumptions

	{
		releaseIndex, err := c.AppOpts.GetReleaseIndex("default")
		if err != nil {
			return errors.Wrap(err, "loading release index")
		}

		stemcellIndex, err := c.AppOpts.GetStemcellIndex("default")
		if err != nil {
			return errors.Wrap(err, "loading stemcell index")
		}

		var rootQuery = graphql.NewObject(
			graphql.ObjectConfig{
				Name: "Query",
				Fields: graphql.Fields{
					"releases":  releaseversiongraphql.NewListQuery(releaseIndex),
					"stemcells": stemcellversiongraphql.NewListQuery(stemcellIndex),
				},
			},
		)
		var schema, _ = graphql.NewSchema(
			graphql.SchemaConfig{
				Query: rootQuery,
			},
		)

		r.HandleFunc("/api/v2/graphql", func(w http.ResponseWriter, r *http.Request) {
			// TODO switch to post?
			result := graphql.Do(graphql.Params{
				Schema:        schema,
				RequestString: r.URL.Query().Get("query"),
			})

			if len(result.Errors) > 0 {
				fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
			}

			w.Header().Set("Access-Control-Allow-Origin", "*")

			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "  ")
			encoder.Encode(result)
		})
	}

	// {
	// 	datastore, err := c.AppOpts.GetStemcellIndex("default")
	// 	if err != nil {
	// 		return errors.Wrap(err, "loading stemcell index")
	// 	}
	//
	// 	stemcellversionv2.NewHandler(c.AppOpts.GetLogger(), datastore, scheduler).RegisterHandlers(r)
	// }

	return http.ListenAndServe(cfg.Bind, r)
}
