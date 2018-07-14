package urlutil

import (
	"net/http"

	"github.com/dpb587/boshua/server/httputil"
	"github.com/dpb587/boshua/stemcellversion"
)

func ApplyStemcellVersionRefToQuery(r *http.Request, ref stemcellversion.Reference) {
	q := r.URL.Query()

	q.Add("stemcell.iaas", ref.IaaS)
	q.Add("stemcell.hypervisor", ref.Hypervisor)
	q.Add("stemcell.os", ref.OS)
	q.Add("stemcell.version", ref.Version)

	r.URL.RawQuery = q.Encode()
}

func StemcellVersionRefFromParam(r *http.Request) (stemcellversion.Reference, error) {
	iaas, err := httputil.SimpleQueryLookup(r, "stemcell.iaas")
	if err != nil {
		return stemcellversion.Reference{}, err
	}

	hypervisor, err := httputil.SimpleQueryLookup(r, "stemcell.hypervisor")
	if err != nil {
		return stemcellversion.Reference{}, err
	}

	os, err := httputil.SimpleQueryLookup(r, "stemcell.os")
	if err != nil {
		return stemcellversion.Reference{}, err
	}

	version, err := httputil.SimpleQueryLookup(r, "stemcell.version")
	if err != nil {
		return stemcellversion.Reference{}, err
	}

	return stemcellversion.Reference{
		IaaS:       iaas,
		Hypervisor: hypervisor,
		OS:         os,
		Version:    version,
	}, nil
}
