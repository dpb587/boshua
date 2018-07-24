package datastore

import (
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
)

type FilterParams struct {
	Release *releaseversiondatastore.FilterParams
	OS      *osversiondatastore.FilterParams
}
