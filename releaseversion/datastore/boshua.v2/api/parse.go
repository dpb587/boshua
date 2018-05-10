package api

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
		return releaseversion.Artifact{}, baseLogger, errors.Wrap(err, "parsing release version")
	}

	logger := baseLogger.WithFields(logrus.Fields{
		"boshua.release.name":    releaseVersionRef.Name,
		"boshua.release.version": releaseVersionRef.Version,
	})

	if len(releaseVersionRef.Checksums) > 0 {
		logger = logger.WithField("boshua.release.checksum", releaseVersionRef.Checksums[0].String())
	}

	releaseVersions, err := ds.Filter(releaseVersionRef)
	if err != nil {
		return releaseversion.Artifact{}, logger, httputil.NewError(err, http.StatusInternalServerError, "release version index failed")
	} else if len(releaseVersions) == 0 {
		return releaseversion.Artifact{}, logger, httputil.NewError(err, http.StatusNotFound, "release version not found")
	} else if len(releaseVersions) > 1 {
		return releaseversion.Artifact{}, logger, httputil.NewError(err, http.StatusBadRequest, "multiple release versions found")
	}

	return releaseVersions[0], logger, nil
}
