package datastore

import (
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
)

type Factory interface {
	Create(provider, name string, options map[string]interface{}, stemcellVersionIndex stemcellversiondatastore.Index) (Index, error)
}
