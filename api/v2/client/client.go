package client

import (
	"net/http"
	"reflect"
	"time"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/sirupsen/logrus"
)

type Client struct {
	client   *http.Client
	endpoint string
	logger   logrus.FieldLogger
}

func New(client *http.Client, endpoint string, logger logrus.FieldLogger) *Client {
	return &Client{
		client:   client,
		endpoint: endpoint,
		logger: logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(Client{}).PkgPath(),
			"api.version":   "v2",
		}),
	}
}

func (c *Client) doRequest(logger logrus.FieldLogger, request *http.Request) (*http.Response, error) {
	var uuidString string

	uuidResult, err := uuid.NewV4()
	if err != nil {
		// TODO bail?
		uuidString = "unknown"
	} else {
		uuidString = uuidResult.String()
	}

	t0 := time.Now()

	httplogger := logger.WithFields(logrus.Fields{
		"http.request.time":   t0.Format(time.RFC3339),
		"http.request.id":     uuidString,
		"http.request.method": request.Method,
		"http.request.uri":    request.URL.String(),
	})

	httplogger.Debugf("sending request")

	response, err := c.client.Do(request)
	if err != nil {
		httplogger.WithField("error", err.Error()).Warnf("errored during request")

		return response, err
	}

	t1 := time.Now()

	httplogger.WithFields(logrus.Fields{
		"http.response.time":     t1.Format(time.RFC3339),
		"http.response.duration": t1.Sub(t0) / time.Millisecond,
		"http.response.status":   response.StatusCode,
	}).Debugf("received response")

	return response, err
}
