package compiledreleaseversion

import (
	"fmt"
	"net/http"

	"github.com/dpb587/boshua/api/v2/urlutil"
	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/sirupsen/logrus"
)

func parseRequest(logger logrus.FieldLogger, r *http.Request) (compiledreleaseversion.Reference, logrus.FieldLogger, error) {
	releaseVersionRef, err := urlutil.ReleaseVersionRefFromParam(r)
	if err != nil {
		return compiledreleaseversion.Reference{}, nil, fmt.Errorf("parsing release version: %v", err)
	}

	osVersionRef, err := urlutil.OSVersionRefFromParam(r)
	if err != nil {
		return compiledreleaseversion.Reference{}, nil, fmt.Errorf("parsing os version: %v", err)
	}

	logger = logger.WithFields(logrus.Fields{
		"boshua.release.name":     releaseVersionRef.Name,
		"boshua.release.version":  releaseVersionRef.Version,
		"boshua.release.checksum": releaseVersionRef.Checksums[0].String(),
		"boshua.os.name":          osVersionRef.Name,
		"boshua.os.version":       osVersionRef.Version,
	})

	ref := compiledreleaseversion.Reference{
		ReleaseVersion: releaseVersionRef,
		OSVersion:      osVersionRef,
	}

	return ref, logger, nil
}
