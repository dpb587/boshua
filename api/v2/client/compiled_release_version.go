package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	api "github.com/dpb587/boshua/api/v2/models/compiledreleaseversion"
	"github.com/dpb587/boshua/api/v2/urlutil"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
)

func (c *Client) GetCompiledReleaseVersionCompilation(releaseVersion releaseversion.Reference, osVersion osversion.Reference) (*api.GETCompilationResponse, error) {
	logger := c.logger.WithField("api.handler", "compiledreleaseversion/info")

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%sv2/compiled-release-version/compilation", c.endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	urlutil.ApplyReleaseVersionRefToQuery(request, releaseVersion)
	urlutil.ApplyOSVersionRefToQuery(request, osVersion)

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

	var res *api.GETCompilationResponse

	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %v", err)
	}

	return res, nil
}

func (c *Client) RequestCompiledReleaseVersionCompilation(releaseVersion releaseversion.Reference, osVersion osversion.Reference) (*api.POSTCompilationResponse, error) {
	logger := c.logger.WithField("api.handler", "compiledreleaseversion/compilation")

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%sv2/compiled-release-version/compilation", c.endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	urlutil.ApplyReleaseVersionRefToQuery(request, releaseVersion)
	urlutil.ApplyOSVersionRefToQuery(request, osVersion)

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

	var res *api.POSTCompilationResponse

	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response body: %v", err)
	}

	return res, nil
}

func (c *Client) RequireCompiledReleaseVersionCompilation(releaseVersion releaseversion.Reference, osVersion osversion.Reference) (*api.GETCompilationResponse, error) {
	resInfo, err := c.GetCompiledReleaseVersionCompilation(releaseVersion, osVersion)
	if err != nil {
		return nil, fmt.Errorf("finding compiled release: %v", err)
	} else if resInfo == nil {
		priorStatus := "unknown"

		for {
			res, err := c.RequestCompiledReleaseVersionCompilation(releaseVersion, osVersion)
			if err != nil {
				return nil, fmt.Errorf("requesting compiled release: %v", err)
			} else if res == nil {
				return nil, fmt.Errorf("unsupported compilation")
			}

			if res.Status != priorStatus {
				fmt.Fprintf(os.Stderr, "compilation status: %s\n", res.Status) // TODO
				priorStatus = res.Status
			}

			if res.Complete {
				break
			}

			time.Sleep(10 * time.Second)
		}

		resInfo, err = c.GetCompiledReleaseVersionCompilation(releaseVersion, osVersion)
		if err != nil {
			return nil, fmt.Errorf("finding compiled release: %v", err)
		} else if resInfo == nil {
			return nil, fmt.Errorf("finding compiled release: unable to fetch expected compilation")
		}
	}

	return resInfo, nil
}
