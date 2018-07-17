package opts

import (
	"fmt"
	"strings"

	"github.com/dpb587/boshua/stemcellversion/datastore/boshioindex"
)

type Stemcell struct {
	Name       string
	IaaS       string
	Hypervisor string
	OS         string
	Version    string
}

func (r Stemcell) String() string {
	return fmt.Sprintf("%s/%s", r.Name, r.Version)
}

func (r *Stemcell) UnmarshalFlag(data string) error {
	// TODO parse better
	split := strings.SplitN(data, "/", -1)
	value := fmt.Sprintf("bosh-stemcell-%s-%s-go_agent", split[1], strings.TrimPrefix(strings.TrimSuffix(split[0], "-go_agent"), "bosh-"))
	parsed := boshioindex.ConvertFileNameToReference(value)
	if parsed == nil {
		return fmt.Errorf("unable to parse stemcell: %s", value)
	}

	r.IaaS = parsed.IaaS
	r.Hypervisor = parsed.Hypervisor
	r.OS = parsed.OS
	r.Version = parsed.Version

	return nil
}
