package setter

import "github.com/dpb587/boshua/config/provider"

type Setter interface {
	SetConfig(*provider.Config)
}
