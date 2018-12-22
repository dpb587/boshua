package manifest

import (
	"fmt"
	"strings"

	"github.com/cppforlife/go-patch/patch"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
)

type StemcellPatch struct {
	Name    string `yaml:"name,omitempty"`
	OS      string `yaml:"os,omitempty"`
	Version string `yaml:"version"`

	pointer patch.Pointer
}

func (r StemcellPatch) Slug() string {
	if r.Name != "" {
		return fmt.Sprintf("%s/%s", r.Name, r.Version)
	}

	return fmt.Sprintf("%s/%s", r.OS, r.Version)
}

func (r StemcellPatch) FilterParams() stemcellversiondatastore.FilterParams {
	if r.Name != "" {
		return stemcellversiondatastore.FilterParamsFromSlug(strings.Join([]string{r.Name, r.Version}, "/"))
	}

	return stemcellversiondatastore.FilterParams{
		OSExpected:      true,
		OS:              r.OS,
		VersionExpected: true,
		Version:         r.Version,
	}
}
