package boshuaV2

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/machinebox/graphql"
	"github.com/sirupsen/logrus"
)

type Client struct {
	client *graphql.Client
	config BoshuaConfig
	logger logrus.FieldLogger
}

func NewClient(client *http.Client, config BoshuaConfig, logger logrus.FieldLogger) *Client {
	return &Client{
		client: graphql.NewClient(fmt.Sprintf("%s/api/v2/graphql", config.URL), graphql.WithHTTPClient(client)),
		config: config,
		logger: logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(Client{}).PkgPath(),
			"api.version":   "v2",
		}),
	}
}

func (c *Client) Execute(req *graphql.Request, responseData interface{}) error {
	ctx := context.Background()

	if err := c.client.Run(ctx, req, &responseData); err != nil {
		return err
	}

	return nil
}

// TODO resurrect logger/fields
// func (c *Client) doRequest(request *http.Request) (*http.Response, error) {
// 	var uuidString string
//
// 	uuidResult, err := uuid.NewV4()
// 	if err != nil {
// 		// TODO bail?
// 		uuidString = "unknown"
// 	} else {
// 		uuidString = uuidResult.String()
// 	}
//
// 	t0 := time.Now()
//
// 	httplogger := c.logger.WithFields(logrus.Fields{
// 		"http.request.time":   t0.Format(time.RFC3339),
// 		"http.request.id":     uuidString,
// 		"http.request.method": request.Method,
// 		"http.request.uri":    request.URL.String(),
// 	})
//
// 	httplogger.Debugf("sending request")
//
// 	response, err := c.client.Do(request)
// 	if err != nil {
// 		httplogger.WithField("error", err.Error()).Warnf("errored during request")
//
// 		return response, err
// 	}
//
// 	t1 := time.Now()
//
// 	httplogger.WithFields(logrus.Fields{
// 		"http.response.time":     t1.Format(time.RFC3339),
// 		"http.response.duration": t1.Sub(t0) / time.Millisecond,
// 		"http.response.status":   response.StatusCode,
// 	}).Debugf("received response")
//
// 	return response, err
// }
