package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	releaseversiongraphql "github.com/dpb587/boshua/releaseversion/graphql"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversiongraphql "github.com/dpb587/boshua/stemcellversion/graphql"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
	schedulerboshuaV2 "github.com/dpb587/boshua/task/scheduler/boshua.v2/graphql"
	// releaseversionv2 "github.com/dpb587/boshua/releaseversion/api/v2/server"
	// stemcellversionv2 "github.com/dpb587/boshua/stemcellversion/api/v2/server"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
)

type GraphqlV2 struct {
	releaseIndex            releaseversiondatastore.Index
	releaseCompilationIndex compilationdatastore.Index
	stemcellIndex           stemcellversiondatastore.Index
	scheduler               schedulerpkg.Scheduler
}

func NewGraphqlV2(releaseIndex releaseversiondatastore.Index, releaseCompilationIndex compilationdatastore.Index, stemcellIndex stemcellversiondatastore.Index, scheduler schedulerpkg.Scheduler) *GraphqlV2 {
	return &GraphqlV2{
		releaseIndex:            releaseIndex,
		releaseCompilationIndex: releaseCompilationIndex,
		stemcellIndex:           stemcellIndex,
		scheduler:               scheduler,
	}
}

func (h *GraphqlV2) Mount(m *mux.Router) {
	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"release":        releaseversiongraphql.NewQuery(h.releaseIndex, h.releaseCompilationIndex),
				"releases":       releaseversiongraphql.NewListQuery(h.releaseIndex),
				"release_labels": releaseversiongraphql.NewLabelsQuery(h.releaseIndex),
				"stemcell":       stemcellversiongraphql.NewQuery(h.stemcellIndex),
				"stemcells":      stemcellversiongraphql.NewListQuery(h.stemcellIndex),
			},
		},
	)
	var mutationType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"scheduleReleaseAnalysis":  schedulerboshuaV2.NewReleaseAnalysisField(h.scheduler, h.releaseIndex),
			"scheduleStemcellAnalysis": schedulerboshuaV2.NewStemcellAnalysisField(h.scheduler, h.stemcellIndex),
		},
	})

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)

	if err != nil {
		panic(err)
	}

	m.HandleFunc("/api/v2/graphql", func(w http.ResponseWriter, r *http.Request) {
		requestBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// TODO !panic
			panic(err)
		}

		// fmt.Printf("< %s\n", strings.TrimSpace(string(requestBytes)))

		var requestBodyObj struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}

		err = json.Unmarshal(requestBytes, &requestBodyObj)
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

		// responseBytes, err := json.MarshalIndent(result, "", "  ")
		responseBytes, err := json.Marshal(result)
		if err != nil {
			panic(err) // TODO !panic
		}

		// fmt.Printf("> %s\n", responseBytes)

		w.Write(responseBytes)
		w.Write([]byte("\n"))
		// encoder := json.NewEncoder(w)
		// encoder.SetIndent("", "  ")
		// encoder.Encode(result)
	}).Methods(http.MethodPost)
}
