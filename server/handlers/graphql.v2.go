package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	releaseversiongraphql "github.com/dpb587/boshua/releaseversion/graphql"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversiongraphql "github.com/dpb587/boshua/stemcellversion/graphql"
	// releaseversionv2 "github.com/dpb587/boshua/releaseversion/api/v2/server"
	// stemcellversionv2 "github.com/dpb587/boshua/stemcellversion/api/v2/server"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
)

type GraphqlV2 struct {
	releaseIndex  releaseversiondatastore.Index
	stemcellIndex stemcellversiondatastore.Index
}

func NewGraphqlV2(releaseIndex releaseversiondatastore.Index, stemcellIndex stemcellversiondatastore.Index) *GraphqlV2 {
	return &GraphqlV2{
		releaseIndex:  releaseIndex,
		stemcellIndex: stemcellIndex,
	}
}

func (h *GraphqlV2) Mount(m *mux.Router) {
	var rootQuery = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"releases":       releaseversiongraphql.NewListQuery(h.releaseIndex),
				"release_labels": releaseversiongraphql.NewLabelsQuery(h.releaseIndex),
				"stemcells":      stemcellversiongraphql.NewListQuery(h.stemcellIndex),
				"stemcell":       stemcellversiongraphql.NewQuery(h.stemcellIndex),
			},
		},
	)
	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: rootQuery,
		},
	)

	m.HandleFunc("/api/v2/graphql", func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// TODO !panic
			panic(err)
		}

		var requestBodyObj struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}

		err = json.Unmarshal(bodyBytes, &requestBodyObj)
		if err != nil {
			// TODO !panic
			panic(err)
		}

		// TODO switch to post?
		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  requestBodyObj.Query,
			VariableValues: requestBodyObj.Variables,
		})

		if len(result.Errors) > 0 {
			fmt.Printf("wrong result, unexpected errors: %v\n", result.Errors)
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(result)
	}).Methods(http.MethodPost)
}
