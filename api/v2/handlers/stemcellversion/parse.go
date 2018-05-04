package stemcellversion

import (
	"fmt"
	"net/http"

	"github.com/dpb587/boshua/api/v2/httputil"
	"github.com/dpb587/boshua/api/v2/urlutil"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/sirupsen/logrus"
)

func parseRequest(baseLogger logrus.FieldLogger, r *http.Request, ds datastore.Index) (stemcellversion.Artifact, logrus.FieldLogger, error) {
	stemcellVersionRef, err := urlutil.StemcellVersionRefFromParam(r)
	if err != nil {
		return stemcellversion.Artifact{}, baseLogger, fmt.Errorf("parsing stemcell version: %v", err)
	}

	logger := baseLogger.WithFields(logrus.Fields{
		"boshua.stemcell.name":    stemcellVersionRef.Name,
		"boshua.stemcell.version": stemcellVersionRef.Version,
	})

	stemcellVersion, err := ds.Find(stemcellVersionRef)
	if err != nil {
		httperr := httputil.NewError(err, http.StatusInternalServerError, "stemcell version index failed")

		if err == datastore.MissingErr {
			httperr = httputil.NewError(err, http.StatusNotFound, "stemcell version not found")
		}

		return stemcellversion.Artifact{}, logger, httperr
	}

	return stemcellVersion, logger, nil
}
