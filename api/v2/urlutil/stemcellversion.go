package urlutil

import (
	"net/http"

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
	iaas, err := simpleQueryLookup(r, "stemcell.iaas")
	if err != nil {
		return stemcellversion.Reference{}, err
	}

	hypervisor, err := simpleQueryLookup(r, "stemcell.hypervisor")
	if err != nil {
		return stemcellversion.Reference{}, err
	}

	os, err := simpleQueryLookup(r, "stemcell.os")
	if err != nil {
		return stemcellversion.Reference{}, err
	}

	version, err := simpleQueryLookup(r, "stemcell.version")
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
