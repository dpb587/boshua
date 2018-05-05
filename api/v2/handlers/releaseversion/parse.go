package releaseversion

import (
	"fmt"
	"net/http"

	"github.com/dpb587/boshua/api/v2/httputil"
	"github.com/dpb587/boshua/api/v2/urlutil"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/sirupsen/logrus"
)

func parseRequest(baseLogger logrus.FieldLogger, r *http.Request, ds datastore.Index) (releaseversion.Artifact, logrus.FieldLogger, error) {
	releaseVersionRef, err := urlutil.ReleaseVersionRefFromParam(r)
	if err != nil {
		return releaseversion.Artifact{}, baseLogger, fmt.Errorf("parsing release version: %v", err)
	}

	logger := baseLogger.WithFields(logrus.Fields{
		"boshua.release.name":    releaseVersionRef.Name,
		"boshua.release.version": releaseVersionRef.Version,
	})

	if len(releaseVersionRef.Checksums) > 0 {
		logger = logger.WithField("boshua.release.checksum", releaseVersionRef.Checksums[0].String())
	}

	releaseVersion, err := ds.Find(releaseVersionRef)
	if err != nil {
		httperr := httputil.NewError(err, http.StatusInternalServerError, "release version index failed")

		if err == datastore.MissingErr {
			httperr = httputil.NewError(err, http.StatusNotFound, "release version not found")
		}

		return releaseversion.Artifact{}, logger, httperr
	}

	return releaseVersion, logger, nil
}
