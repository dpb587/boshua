package stemcellversion

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/dpb587/boshua"
)

type Subject struct {
	Reference

	MetalinkSource map[string]interface{}
}

var _ boshua.Subject = &Subject{}

func (s Subject) SubjectReference() boshua.Reference {
	return boshua.Reference{
		Context: "stemcellversion",
		ID:      s.id(),
	}
}

func (s Subject) id() string {
	h := sha1.New()
	h.Write([]byte(strings.Join(
		[]string{
			s.Reference.OS,
			s.Reference.Version,
		},
		"/",
	)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
