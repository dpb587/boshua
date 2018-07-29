package params

import (
	"errors"
	"net/http"

	"github.com/dpb587/boshua/stemcellversion/datastore"
)

func FilterParamsFromQuery(r *http.Request) (datastore.FilterParams, error) {
	q := r.URL.Query()

	f := datastore.FilterParams{}

	v, ok := q["stemcell-os"]
	if ok {
		if len(v) != 1 {
			return datastore.FilterParams{}, errors.New("stemcell-os: expected single value")
		}

		f.OSExpected = true
		f.OS = v[0]
	}

	v, ok = q["stemcell-version"]
	if ok {
		if len(v) != 1 {
			return datastore.FilterParams{}, errors.New("stemcell-version: expected single value")
		}

		f.VersionExpected = true
		f.Version = v[0]
	}

	v, ok = q["stemcell-iaas"]
	if ok {
		if len(v) != 1 {
			return datastore.FilterParams{}, errors.New("stemcell-iaas: expected single value")
		}

		f.IaaSExpected = true
		f.IaaS = v[0]
	}

	v, ok = q["stemcell-hypervisor"]
	if ok {
		if len(v) != 1 {
			return datastore.FilterParams{}, errors.New("stemcell-hypervisor: expected single value")
		}

		f.HypervisorExpected = true
		f.Hypervisor = v[0]
	}

	v, ok = q["stemcell-flavor"]
	if ok {
		if len(v) != 1 {
			return datastore.FilterParams{}, errors.New("stemcell-flavor: expected single value")
		}

		f.FlavorExpected = true
		f.Flavor = v[0]
	}

	v, ok = q["stemcell-label"]
	if ok && len(v) > 0 {
		f.LabelsExpected = true
		f.Labels = v
	}

	return f, nil
}
