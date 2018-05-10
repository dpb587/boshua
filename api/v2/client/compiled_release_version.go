package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	analysisapi "github.com/dpb587/boshua/api/v2/models/analysis"
	api "github.com/dpb587/boshua/api/v2/models/compiledreleaseversion"
	schedulerapi "github.com/dpb587/boshua/api/v2/models/scheduler"
	"github.com/dpb587/boshua/api/v2/urlutil"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
)

func (c *Client) GetCompiledReleaseVersionCompilation(releaseVersion releaseversion.Reference, osVersion osversion.Reference) (*api.GETCompilationInfoResponse, error) {
	logger := c.logger.WithField("api.handler", "compiledreleaseversion/info")

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%sv2/compiled-release-version/compilation/info", c.endpoint), nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	urlutil.ApplyReleaseVersionRefToQuery(request, releaseVersion)
	urlutil.ApplyOSVersionRefToQuery(request, osVersion)

	response, err := c.doRequest(logger, request)
	if err != nil {
		return nil, errors.Wrap(err, "executing request")
	} else if response.StatusCode == http.StatusNotFound {
		// not available; expected
		return nil, nil
	} else if response.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("executing request: status %d: %s", response.StatusCode, bodyBytes)
	}

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading response body")
	}

	var res *api.GETCompilationInfoResponse

	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling response body")
	}

	return res, nil
}

func (c *Client) RequestCompiledReleaseVersionCompilation(releaseVersion releaseversion.Reference, osVersion osversion.Reference) (*api.POSTCompilationQueueResponse, error) {
	logger := c.logger.WithField("api.handler", "compiledreleaseversion/compilation")

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%sv2/compiled-release-version/compilation/queue", c.endpoint), nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	urlutil.ApplyReleaseVersionRefToQuery(request, releaseVersion)
	urlutil.ApplyOSVersionRefToQuery(request, osVersion)

	response, err := c.doRequest(logger, request)
	if err != nil {
		return nil, errors.Wrap(err, "executing request")
	} else if response.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("executing request: status %d: %s", response.StatusCode, bodyBytes)
	}

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading response body")
	}

	var res *api.POSTCompilationQueueResponse

	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling response body")
	}

	return res, nil
}

func (c *Client) RequireCompiledReleaseVersionCompilation(releaseVersion releaseversion.Reference, osVersion osversion.Reference, taskStatusWatcher TaskStatusWatcher) (*api.GETCompilationInfoResponse, error) {
	resInfo, err := c.GetCompiledReleaseVersionCompilation(releaseVersion, osVersion)
	if err != nil {
		return nil, errors.Wrap(err, "finding compiled release")
	} else if resInfo == nil {
		priorStatus := schedulerapi.TaskStatus{}

		for {
			res, err := c.RequestCompiledReleaseVersionCompilation(releaseVersion, osVersion)
			if err != nil {
				return nil, errors.Wrap(err, "requesting compiled release")
			} else if res == nil {
				return nil, fmt.Errorf("unsupported compilation")
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

		resInfo, err = c.GetCompiledReleaseVersionCompilation(releaseVersion, osVersion)
		if err != nil {
			return nil, errors.Wrap(err, "finding compiled release")
		} else if resInfo == nil {
			return nil, fmt.Errorf("finding compiled release: unable to fetch expected compilation")
		}
	}

	return resInfo, nil
}

func (c *Client) GetCompiledReleaseVersionAnalysis(releaseVersion releaseversion.Reference, osVersion osversion.Reference, analyzer string) (*analysisapi.GETInfoResponse, error) {
	logger := c.logger.WithField("api.handler", "compiledreleaseversion/analysis")

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%sv2/compiled-release-version/analysis/info", c.endpoint), nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	urlutil.ApplyReleaseVersionRefToQuery(request, releaseVersion)
	urlutil.ApplyOSVersionRefToQuery(request, osVersion)
	urlutil.ApplyAnalysisAnalyzerToQuery(request, analyzer)

	response, err := c.doRequest(logger, request)
	if err != nil {
		return nil, errors.Wrap(err, "executing request")
	} else if response.StatusCode == http.StatusNotFound {
		// not available; expected
		return nil, nil
	} else if response.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("executing request: status %d: %s", response.StatusCode, bodyBytes)
	}

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading response body")
	}

	var res *analysisapi.GETInfoResponse

	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling response body")
	}

	return res, nil
}

func (c *Client) RequestCompiledReleaseVersionAnalysis(releaseVersion releaseversion.Reference, osVersion osversion.Reference, analyzer string) (*analysisapi.POSTQueueResponse, error) {
	logger := c.logger.WithField("api.handler", "compiledreleaseversion/analysis")

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%sv2/compiled-release-version/analysis/queue", c.endpoint), nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	urlutil.ApplyReleaseVersionRefToQuery(request, releaseVersion)
	urlutil.ApplyOSVersionRefToQuery(request, osVersion)
	urlutil.ApplyAnalysisAnalyzerToQuery(request, analyzer)

	response, err := c.doRequest(logger, request)
	if err != nil {
		return nil, errors.Wrap(err, "executing request")
	} else if response.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("executing request: status %d: %s", response.StatusCode, bodyBytes)
	}

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading response body")
	}

	var res *analysisapi.POSTQueueResponse

	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling response body")
	}

	return res, nil
}

func (c *Client) RequireCompiledReleaseVersionAnalysis(releaseVersion releaseversion.Reference, osVersion osversion.Reference, analyzer string, taskStatusWatcher TaskStatusWatcher) (*analysisapi.GETInfoResponse, error) {
	resInfo, err := c.GetCompiledReleaseVersionAnalysis(releaseVersion, osVersion, analyzer)
	if err != nil {
		return nil, errors.Wrap(err, "finding analysis")
	} else if resInfo == nil {
		priorStatus := schedulerapi.TaskStatus{}

		for {
			res, err := c.RequestCompiledReleaseVersionAnalysis(releaseVersion, osVersion, analyzer)
			if err != nil {
				return nil, errors.Wrap(err, "requesting analysis")
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

		resInfo, err = c.GetCompiledReleaseVersionAnalysis(releaseVersion, osVersion, analyzer)
		if err != nil {
			return nil, errors.Wrap(err, "finding analysis")
		} else if resInfo == nil {
			return nil, fmt.Errorf("finding analysis: unable to fetch expected analysis")
		}
	}

	return resInfo, nil
}
