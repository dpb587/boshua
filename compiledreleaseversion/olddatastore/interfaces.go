package datastore

import "github.com/dpb587/boshua/analysis/datastore"

type Factory interface {
	Create(provider, name string, options map[string]interface{}, analysisIndex datastore.Index) (Index, error)
}
