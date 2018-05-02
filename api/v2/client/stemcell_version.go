package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	api "github.com/dpb587/boshua/api/v2/models/analysis"
	schedulerapi "github.com/dpb587/boshua/api/v2/models/scheduler"
	"github.com/dpb587/boshua/api/v2/urlutil"
	"github.com/dpb587/boshua/stemcellversion"
)

func (c *Client) GetStemcellVersionAnalysis(stemcellVersion stemcellversion.Reference, analyzer string) (*api.GETAnalysisResponse, error) {
	logger := c.logger.WithField("api.handler", "stemcellversion/analysis")

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%sv2/stemcell-version/analysis/info", c.endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	urlutil.ApplyStemcellVersionRefToQuery(request, stemcellVersion)
	urlutil.ApplyAnalysisAnalyzerToQuery(request, analyzer)

	response, err := c.doRequest(logger, request)
	if err != nil {
		return nil, fmt.Errorf("executing request: %v", err)
	} else if response.StatusCode == http.StatusNotFound {
		// not available; expected
		return nil, nil
	} else if response.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("executing request: status %d: %s", response.StatusCode, bodyBytes)
	}

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	var res *api.GETAnalysisResponse

	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %v", err)
	}

	return res, nil
}

func (c *Client) RequestStemcellVersionAnalysis(stemcellVersion stemcellversion.Reference, analyzer string) (*api.POSTAnalysisResponse, error) {
	logger := c.logger.WithField("api.handler", "stemcellversion/analysis")

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%sv2/stemcell-version/analysis/queue", c.endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	urlutil.ApplyStemcellVersionRefToQuery(request, stemcellVersion)
	urlutil.ApplyAnalysisAnalyzerToQuery(request, analyzer)

	response, err := c.doRequest(logger, request)
	if err != nil {
		return nil, fmt.Errorf("executing request: %v", err)
	} else if response.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("executing request: status %d: %s", response.StatusCode, bodyBytes)
	}

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	var res *api.POSTAnalysisResponse

	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %v", err)
	}

	return res, nil
}

func (c *Client) RequireStemcellVersionAnalysis(stemcellVersion stemcellversion.Reference, analyzer string, taskStatusWatcher TaskStatusWatcher) (*api.GETAnalysisResponse, error) {
	resInfo, err := c.GetStemcellVersionAnalysis(stemcellVersion, analyzer)
	if err != nil {
		return nil, fmt.Errorf("finding analysis: %v", err)
	} else if resInfo == nil {
		priorStatus := schedulerapi.TaskStatus{}

		for {
			res, err := c.RequestStemcellVersionAnalysis(stemcellVersion, analyzer)
			if err != nil {
				return nil, fmt.Errorf("requesting analysis: %v", err)
			} else if res == nil {
				return nil, fmt.Errorf("unsupported analysis")
			}

			currentStatus := res.Data

			if currentStatus != priorStatus {
				if taskStatusWatcher != nil {
					taskStatusWatcher(currentStatus)
				}

				priorStatus = currentStatus
			}

			if currentStatus.Complete {
				break
			}

			time.Sleep(10 * time.Second)
		}

		resInfo, err = c.GetStemcellVersionAnalysis(stemcellVersion, analyzer)
		if err != nil {
			return nil, fmt.Errorf("finding analysis: %v", err)
		} else if resInfo == nil {
			return nil, fmt.Errorf("finding analysis: unable to fetch expected analysis")
		}
	}

	return resInfo, nil
}