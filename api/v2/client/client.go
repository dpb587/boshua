package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dpb587/bosh-compiled-releases/api/v2/models"
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
		logger:   logger,
	}
}

func (c *Client) CompiledReleaseVersionInfo(req models.CRVInfoRequest) (*models.CRVInfoResponse, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshalling request data: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%sv2/compiled-release-version/info", c.endpoint), strings.NewReader(string(reqBytes)))
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	response, err := c.client.Do(request)
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
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshalling request data: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%sv2/compiled-release-version/request", c.endpoint), strings.NewReader(string(reqBytes)))
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	response, err := c.client.Do(request)
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
