package client

import (
	"fmt"
	"os"
	"time"

	"github.com/dpb587/boshua/api/v2/models"
)

func RequestAndWait(client *Client, releaseRef models.ReleaseVersionRef, stemcellRef models.OSVersionRef) (*models.CRVInfoResponse, error) {
	resInfo, err := client.CompiledReleaseVersionInfo(models.CRVInfoRequest{
		Data: models.CRVInfoRequestData{
			ReleaseVersionRef:  releaseRef,
			OSVersionRef: stemcellRef,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("finding compiled release: %v", err)
	} else if resInfo == nil {
		priorStatus := "unknown"

		for {
			res, err := client.CompiledReleaseVersionRequest(models.CRVRequestRequest{
				Data: models.CRVRequestRequestData{
					ReleaseVersionRef:  releaseRef,
					OSVersionRef: stemcellRef,
				},
			})
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

		resInfo, err = client.CompiledReleaseVersionInfo(models.CRVInfoRequest{
			Data: models.CRVInfoRequestData{
				ReleaseVersionRef:  releaseRef,
				OSVersionRef: stemcellRef,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("finding compiled release: %v", err)
		} else if resInfo == nil {
			return nil, fmt.Errorf("finding compiled release: unable to fetch expected compilation")
		}
	}

	return resInfo, nil
}
