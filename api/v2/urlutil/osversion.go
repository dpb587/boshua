package urlutil

import (
	"net/http"

	"github.com/dpb587/boshua/osversion"
)

func ApplyOSVersionRefToQuery(r *http.Request, ref osversion.Reference) {
	q := r.URL.Query()

	q.Add("os.name", ref.Name)
	q.Add("os.version", ref.Version)

	r.URL.RawQuery = q.Encode()
}

func OSVersionRefFromParam(r *http.Request) (osversion.Reference, error) {
	osName, err := simpleQueryLookup(r, "os.name")
	if err != nil {
		return osversion.Reference{}, err
	}

	osVersion, err := simpleQueryLookup(r, "os.version")
	if err != nil {
		return osversion.Reference{}, err
	}

	return osversion.Reference{
		Name:    osName,
		Version: osVersion,
	}, nil
}
