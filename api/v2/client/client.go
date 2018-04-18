package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
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

func (c *Client) CompiledReleaseVersionInfo(req models.CRVInfoRequest) (*models.CRVInfoResponse, error) {
	logger := c.logger.WithField("api.handler", "compiledreleaseversion/info")

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshalling request data: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%sv2/compiled-release-version/info", c.endpoint), strings.NewReader(string(reqBytes)))
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	response, err := c.doRequest(logger, request)
	if err != nil {
		return nil, fmt.Errorf("executing request: %v", err)
	} else if response.StatusCode == http.StatusNotFound {
		// not available; expected
		return nil, nil
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("executing request: status %d", response.StatusCode)
	}

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	var res *models.CRVInfoResponse

	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %v", err)
	}

	return res, nil
}

func (c *Client) CompiledReleaseVersionRequest(req models.CRVRequestRequest) (*models.CRVRequestResponse, error) {
	logger := c.logger.WithField("api.handler", "compiledreleaseversion/request")

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshalling request data: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%sv2/compiled-release-version/request", c.endpoint), strings.NewReader(string(reqBytes)))
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	response, err := c.doRequest(logger, request)
	if err != nil {
		return nil, fmt.Errorf("executing request: %v", err)
	} else if response.StatusCode == http.StatusNotFound {
		// not available; expected
		return nil, nil
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("executing request: status %d", response.StatusCode)
	}

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	var res *models.CRVRequestResponse

	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %v", err)
	}

	return res, nil
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
