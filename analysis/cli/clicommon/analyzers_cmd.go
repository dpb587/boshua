package clicommon

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type AnalyzersCmd struct{}

func (c *AnalyzersCmd) Execute(loader SubjectLoader) error {
	subject, err := loader()
	if err != nil {
		return errors.Wrap(err, "finding subject")
	}

	fmt.Printf("%s\n", strings.Join(subject.SupportedAnalyzers(), "\n"))

	return nil
}
