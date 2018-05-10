package boshuav2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/sirupsen/logrus"
)

type Client struct {
	client *http.Client
	config BoshuaConfig
	logger logrus.FieldLogger
}

func NewClient(client *http.Client, config BoshuaConfig, logger logrus.FieldLogger) *Client {
	return &Client{
		client: client,
		config: config,
		logger: logger.WithFields(logrus.Fields{
			"build.package": reflect.TypeOf(Client{}).PkgPath(),
			"api.version":   "v2",
		}),
	}
}

func (c *Client) Execute(request *http.Request, responseData interface{}) error {
	// TODO prefix endpoint

	response, err := c.doRequest(request)
	if err != nil {
		return fmt.Errorf("executing request: %v", err)
	} else if response.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("executing request: status %d: %s", response.StatusCode, bodyBytes)
	}

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	err = json.Unmarshal(resBytes, responseData)
	if err != nil {
		return fmt.Errorf("unmarshalling response body: %v", err)
	}

	return nil
}

func (c *Client) doRequest(request *http.Request) (*http.Response, error) {
	var uuidString string

	uuidResult, err := uuid.NewV4()
	if err != nil {
		// TODO bail?
		uuidString = "unknown"
	} else {
		uuidString = uuidResult.String()
	}

	t0 := time.Now()

	httplogger := c.logger.WithFields(logrus.Fields{
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
